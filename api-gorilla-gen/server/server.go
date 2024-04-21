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
	"log"
	"net/http"

	chiMiddleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
)

func NewValidator() func(http.Handler) http.Handler {
	spec, err := GetSwagger()
	if err != nil {
		panic(fmt.Errorf("loading spec: %w", err))
	}

	validator := chiMiddleware.OapiRequestValidatorWithOptions(spec, &chiMiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: func(fctx context.Context, ai *openapi3filter.AuthenticationInput) error {
				// ctx := ai.RequestValidationInput.Request.Context()
				// requestUser, _ := ctx.Value("RequestUser").(string)
				// logger := log.Logger.With().
				// 	Str("requestID", middleware.GetReqID(ctx)).
				// 	Str("requestUser", requestUser).
				// 	Logger()
				// if requestUser == "" {
				// 	logger.Info().Msgf("validator: rejected")
				// 	return errors.New("ErrUnauthenticated")
				// } else {
				// 	logger.Info().Msgf("validator: pass")
				// 	return nil
				// }
				return nil
			},
		},
	})
	return validator
}

func Authenticator0(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authenticated := false

		if !authenticated {
			authHdr := r.Header.Get("X-API-Key")
			log.Printf("Try ApiKeyAuth: %v", authHdr)
			if authHdr == "a" {
				authenticated = true
				ctx = context.WithValue(ctx, "AuthenticationBackend", "ApiKeyAuth")
				ctx = context.WithValue(ctx, "RequestUser", authHdr)
				log.Printf("ApiKeyAuth: Passed")
			}
		}
		if authenticated {
			log.Printf("Authenticated")
			r = r.WithContext(ctx)
		} else {
			log.Printf("Anonymous access")
			ctx = context.WithValue(ctx, "AuthenticationBackend", "")
			ctx = context.WithValue(ctx, "RequestUser", "")
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Authenticator(fn http.HandlerFunc) http.HandlerFunc {
	// fn := func(w http.ResponseWriter, r *http.Request) {
	// 	ctx := r.Context()
	// 	authenticated := false

	// 	if !authenticated {
	// 		authHdr := r.Header.Get("X-API-Key")
	// 		log.Printf("Try ApiKeyAuth: %v", authHdr)
	// 		if authHdr == "a" {
	// 			authenticated = true
	// 			ctx = context.WithValue(ctx, "AuthenticationBackend", "ApiKeyAuth")
	// 			ctx = context.WithValue(ctx, "RequestUser", authHdr)
	// 			log.Printf("ApiKeyAuth: Passed")
	// 		}
	// 	}
	// 	if authenticated {
	// 		log.Printf("Authenticated")
	// 		r = r.WithContext(ctx)
	// 	} else {
	// 		log.Printf("Anonymous access")
	// 		ctx = context.WithValue(ctx, "AuthenticationBackend", "")
	// 		ctx = context.WithValue(ctx, "RequestUser", "")
	// 		r = r.WithContext(ctx)
	// 	}
	// 	next.ServeHTTP(w, r)
	// }

	return fn
}

func ShouldBeAuthenticated(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authenticationBackend := ctx.Value("AuthenticationBackend").(string)
		requestUser := ctx.Value("RequestUser").(string)

		if authenticationBackend == "" || requestUser == "" {
			log.Printf("Access: Permision Denied, User=%v", requestUser)
			w.WriteHeader(http.StatusForbidden)
			_, err := io.WriteString(w, "Permision Denied")
			if err != nil {
				log.Printf("Error - %v", err)
			}
			return
		}

		log.Printf("Access: Pass, User=%v", requestUser)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

var (
	ErrUnauthenticated = errors.New("ErrUnauthenticated")
	ErrApiKey          = errors.New("ErrApiKey")
	ErrBasic           = errors.New("ErrBasic")
)

func AuthenticateAuthenticated(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	return nil
}

type ServerImpl struct {
}

func NewServer() ServerInterface {
	return &ServerImpl{}
}

func (s *ServerImpl) GetStatusV1(w http.ResponseWriter, r *http.Request, params GetStatusV1Params) {
	log.Print("GetStatusV1 >")
	isFull := false
	if params.IsFull != nil {
		isFull = *params.IsFull
	}
	status := GetStatusV1{
		Status: "Ok or Not",
	}
	if isFull {
		var cpu int64 = -101
		status.Cpu = &cpu
	}
	log.Printf("GetStatusV1 < %v", status)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(status)
}

func (s *ServerImpl) SetStatusV1(w http.ResponseWriter, r *http.Request) {
	log.Print("GetStatusV1 >")
	var status GetStatusV1

	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		log.Printf("invalid json: %v", err)
		w.WriteHeader(400)
		return
	}

	log.Printf("GetStatusV1 > %v", status)
	status.Status = "Ok or Not"

	log.Printf("GetStatusV1 < %v", status)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(status)
}
