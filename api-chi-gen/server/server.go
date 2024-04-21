// go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest --config=server.cfg.yaml ../api.yaml
// go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.cfg.yaml ../api.yaml
// go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
// go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
//
//go:generate oapi-codegen --config=model.cfg.yaml  ../../api/api.yaml
//go:generate oapi-codegen --config=client.cfg.yaml ../../api/api.yaml
//go:generate oapi-codegen --config=server.cfg.yaml ../../api/api.yaml

package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	chiMiddleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/rs/zerolog/log"
)

func NewValidator() (func(http.Handler) http.Handler, error) {
	spec, err := GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	validator := chiMiddleware.OapiRequestValidatorWithOptions(spec, &chiMiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: func(fctx context.Context, ai *openapi3filter.AuthenticationInput) error {
				ctx := ai.RequestValidationInput.Request.Context()
				requestUser, _ := ctx.Value("RequestUser").(string)
				logger := log.Logger.With().
					Str("requestID", middleware.GetReqID(ctx)).
					Str("requestUser", requestUser).
					Logger()
				if requestUser == "" {
					logger.Info().Msgf("validator: rejected")
					return errors.New("ErrUnauthenticated")
				} else {
					logger.Info().Msgf("validator: pass")
					return nil
				}
			},
		},
	})
	return validator, nil
}

func Authenticator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		logger := log.Logger.With().
			Str("RequestID", middleware.GetReqID(r.Context())).
			Logger()

		ctx := r.Context()
		authenticated := false

		if !authenticated {
			authHdr := r.Header.Get("X-API-Key")
			logger.Info().
				Str("Header", "X-API-Key").
				Str("Value", authHdr).
				Msg("Try ApiKeyAuth")
			if authHdr == "a" {
				authenticated = true
				ctx = context.WithValue(ctx, "AuthenticationBackend", "ApiKeyAuth")
				ctx = context.WithValue(ctx, "RequestUser", authHdr)
				logger := logger.With().Str("RequestUser", authHdr).Logger()
				logger.Info().Msg("Pass ApiKey")
			}
		}
		if !authenticated {
			username, password, ok := r.BasicAuth()
			logger.Info().
				Bool("ok", ok).
				Str("username", username).
				Str("password", password).
				Msg("Try BasicAuth")
			if ok {
				if username == "a" && password == "a" {
					authenticated = true
					ctx = context.WithValue(ctx, "AuthenticationBackend", "BasicAuth")
					ctx = context.WithValue(ctx, "RequestUser", username)
					logger := logger.With().Str("RequestUser", username).Logger()
					logger.Info().Msg("Pass BasicAuth")
				}
			}
		}
		if authenticated {
			logger.Info().Msg("Authenticated")
			r = r.WithContext(ctx)
		} else {
			logger.Info().Msg("Anonymous access")
			ctx = context.WithValue(ctx, "AuthenticationBackend", "")
			ctx = context.WithValue(ctx, "RequestUser", "")
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func ShouldBeAuthenticated(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authenticationBackend := ctx.Value("AuthenticationBackend").(string)
		requestUser := ctx.Value("RequestUser").(string)

		logger := log.Logger.With().
			Str("RequestID", middleware.GetReqID(r.Context())).
			Str("AuthenticationBackend", authenticationBackend).
			Str("RequestUser", requestUser).
			Logger()

		if authenticationBackend == "" || requestUser == "" {
			logger.Info().Msg("Access: Permision Denied")
			w.WriteHeader(http.StatusForbidden)
			_, err := io.WriteString(w, "Permision Denied")
			if err != nil {
				logger.Err(err).Msg("")
			}
			return
		}

		logger.Info().Msg("Access: Pass")
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

var (
	ErrCodeBug         = errors.New("ErrCodeBug")
	ErrUnauthenticated = errors.New("ErrUnauthenticated")
	ErrApiKey          = errors.New("ErrApiKey")
	ErrBasic           = errors.New("ErrBasic")
)

type ServerImpl struct {
}

func NewServer() ServerInterface {
	return &ServerImpl{}
}

func (s *ServerImpl) GetStatusV1(w http.ResponseWriter, r *http.Request, params GetStatusV1Params) {
	logger := log.Logger.With().
		Str("RequestID", middleware.GetReqID(r.Context())).
		Logger()

	logger.Debug().
		// Str("RequestID", c.Get("RequestID").(string)).
		Msg("GetStatusV1")
	isFull := false
	if params.IsFull != nil {
		isFull = *params.IsFull
	}
	q := ""
	if params.Q != nil {
		q = *params.Q
	}
	logger.Info().
		Bool("IsFull", isFull).
		Int32("XPage", *params.XPage).
		Int32("XPageSize", *params.XPageSize).
		Str("Q", q).
		Msg("GetStatusV1")
	status := GetStatusV1{
		Status: "Ok or Not",
	}
	if isFull {
		var cpu int64 = -101
		status.Cpu = &cpu
	}
	logger.Info().
		Interface("result", status).
		Msg("output")
	render.Status(r, http.StatusOK)
	render.JSON(w, r, status)
}

func (s *ServerImpl) SetStatusV1(w http.ResponseWriter, r *http.Request) {
	logger := log.Logger.With().
		Str("RequestID", middleware.GetReqID(r.Context())).
		Logger()

	logger.Info().
		Msg("SetStatusV1")
	var status SetStatusV1JSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		render.Status(r, http.StatusBadRequest)
		logger.Warn().Err(err).Send()
		return
	}

	logger.Info().
		Interface("value", status).
		Msg("input")
	status.Status = "Ok or Not"
	logger.Info().
		Interface("result", status).
		Msg("output")

	render.Status(r, http.StatusOK)
	render.JSON(w, r, status)
}

type StrictServerImpl struct {
}

func NewStrictServer() StrictServerInterface {
	return &StrictServerImpl{}
}

func (s *StrictServerImpl) GetStatusV1(ctx context.Context, request GetStatusV1RequestObject) (GetStatusV1ResponseObject, error) {
	reqId := middleware.GetReqID(ctx)
	logger := log.Logger.With().
		Str("mode", "strict").Logger()
	logger.Debug().
		Interface("reqId2", reqId).
		Interface("input", request).
		Msg("GetStatusV1")
	result := GetStatusV1200JSONResponse{
		Body: GetStatusV1{
			Status: "Ok or Not",
		},
		Headers: GetStatusV1200ResponseHeaders{
			XPage:      1,
			XPageSize:  20,
			XPageCount: 10,
		},
	}
	logger.Info().
		Interface("result", result).
		Msg("output")
	return result, nil
}

func (s *StrictServerImpl) SetStatusV1(ctx context.Context, request SetStatusV1RequestObject) (SetStatusV1ResponseObject, error) {
	userFromContext := ctx.Value("RequestUser")
	logger := log.Logger.With().
		Str("mode", "strict").
		Interface("userFromContext", userFromContext).
		Logger()
	logger.Debug().
		Interface("input", request).
		Msg("SetStatusV1")
	result := SetStatusV1200JSONResponse{
		Status: "Ok or Not",
	}
	logger.Info().
		Interface("result", result).
		Msg("output")
	return result, nil
}
