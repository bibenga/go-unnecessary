// go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest --config=server.cfg.yaml ../api.yaml
// go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.cfg.yaml ../api.yaml
// go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
// go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
//
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=model.cfg.yaml ../../api/api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=client.cfg.yaml ../../api/api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=server.cfg.yaml ../../api/api.yaml

package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	echoMiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetRequestId(c echo.Context, requestId string) {
	c.Set("RequestId", requestId)
}

func GetRequestId(c echo.Context) string {
	requestId, ok := c.Get("RequestId").(string)
	if ok {
		return requestId
	} else {
		return ""
	}
}

func SetUser(c echo.Context, backend, user string) {
	c.Set("AuthenticationBackend", backend)
	c.Set("RequestUser", user)
}

func GetAuthenticationBackend(c echo.Context) string {
	backend, ok := c.Get("AuthenticationBackend").(string)
	if ok {
		return backend
	} else {
		return ""
	}
}

func GetUser(c echo.Context) string {
	user, ok := c.Get("RequestUser").(string)
	if ok {
		return user
	} else {
		return ""
	}
}

type ApiRequestLoggerHook struct {
	ctx echo.Context
}

func (h ApiRequestLoggerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	// e.Str("authenticationBackend", GetAuthenticationBackend(h.ctx))
	e.Str("requestId", GetRequestId(h.ctx))
	e.Str("user", GetUser(h.ctx))
}

func SetApitLogger(ctx echo.Context, parent zerolog.Logger) {
	// logger := log.With().
	// 	Str("requestId", GetRequestId(ctx)).
	// 	Str("user", user).
	// 	Logger().Hook()
	// logger := log.With().
	// 	Logger().Hook(&ApiLoggerHook{ctx: ctx})
	logger := log.With().
		Logger().Hook(&ApiRequestLoggerHook{ctx: ctx})
	// ctx.Set("ApiLoggerParent", parent)
	ctx.Set("ApiLoggerRequest", logger)
}

func GetApiLogger(ctx echo.Context) zerolog.Logger {
	return ctx.Get("ApiLoggerRequest").(zerolog.Logger)
}

func ApiLoggerHandler(root zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			SetApitLogger(ctx, root)
			return next(ctx)
		}
	}
}

type contextKey string

const EchoContextKey = contextKey("unnecessary/EchoContext")

func StoreEchoContextInRequestContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()
			ctx := r.Context()
			ctx = context.WithValue(ctx, EchoContextKey, c)
			c.SetRequest(r.WithContext(ctx))
			return next(c)
		}
	}
}

func Validator() echo.MiddlewareFunc {
	spec, err := GetSwagger()
	if err != nil {
		panic(fmt.Errorf("loading spec: %w", err))
	}
	validator := echoMiddleware.OapiRequestValidatorWithOptions(spec,
		&echoMiddleware.Options{
			Skipper: func(c echo.Context) bool {
				return strings.HasPrefix(c.Request().URL.Path, "/docs")
			},
			Options: openapi3filter.Options{
				// AuthenticationFunc: Authenticate,
				AuthenticationFunc: AuthenticateAuthenticated,
			},
		})

	return validator
}

var (
	ErrUnauthenticated = errors.New("ErrUnauthenticated")
	ErrApiKey          = errors.New("ErrApiKey")
	ErrBasic           = errors.New("ErrBasic")
)

// func Authenticate(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
// 	c := echoMiddleware.GetEchoContext(ctx)
// 	if c == nil {
// 		log.Panic().
// 			Msg("Echo context was not found")
// 		return ErrUnauthenticated
// 	}
// 	logger := GetApiLogger(c)
// 	// logger := log.Logger
// 	logger.Info().
// 		Msg("Try authenticate")
// 	if input != nil {
// 		logger.Info().
// 			Interface("SecuritySchemeName", input.SecuritySchemeName).
// 			Msg("Authenticate")
// 		request := input.RequestValidationInput.Request
// 		if input.SecuritySchemeName == "ApiKeyAuth" {
// 			logger.Info().
// 				Msg("Try ApiKey")
// 			authHdr := request.Header.Get(input.SecurityScheme.Name)
// 			logger.Info().
// 				Str("Header", input.SecurityScheme.Name).
// 				Str("Value", authHdr).
// 				Msg("Try authenenticate user by apiKey")
// 			if authHdr == "a" {
// 				SetUser(c, "ApiKeyAuth", authHdr)
// 				logger.Info().
// 					Msg("Pass ApiKey")
// 				return nil
// 			}
// 			logger.Info().
// 				Msg("Reject ApiKey")
// 			return ErrApiKey
// 		} else if input.SecuritySchemeName == "BasicAuth" {
// 			logger.Info().
// 				Msg("Try Basic")
// 			username, password, ok := request.BasicAuth()
// 			logger.Info().
// 				Bool("ok", ok).
// 				Str("username", username).
// 				Str("password", "******").
// 				Msg("Try authenenticate user by username")
// 			if ok {
// 				if username == "a" && password == "a" {
// 					SetUser(c, "BasicAuth", username)
// 					logger.Info().
// 						Msg("Pass Basic")
// 					return nil
// 				}
// 			}
// 			logger.Info().
// 				Msg("Reject basic")
// 			return ErrBasic
// 		}
// 	}
// 	logger.Info().
// 		Msg("Reject")
// 	return ErrUnauthenticated
// }

func AuthenticateAuthenticated(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	c := echoMiddleware.GetEchoContext(ctx)
	if c == nil {
		log.Panic().
			Msg("Echo context was not found")
		return ErrUnauthenticated
	}
	logger := GetApiLogger(c)
	// logger := log.Logger

	requestUser := GetUser(c)
	if requestUser == "" {
		logger.Info().Msgf("validator: rejected")
		return errors.New("ErrUnauthenticated")
	} else {
		logger.Info().Msgf("validator: pass")
		return nil
	}
}

func Authenticator() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger := GetApiLogger(c)

			authenticated := false

			if !authenticated {
				authHdr := c.Request().Header.Get("X-API-Key")
				logger.Info().
					Str("Header", "X-API-Key").
					Str("Value", authHdr).
					Msg("Try ApiKey")
				if authHdr == "a" {
					authenticated = true
					SetUser(c, "ApiKeyAuth", authHdr)
					logger.Info().
						Msg("Pass ApiKey")
				}
			}

			if !authenticated {
				username, password, ok := c.Request().BasicAuth()
				logger.Info().
					Bool("ok", ok).
					Str("username", username).
					Str("password", password).
					Msg("Try Basic")
				if ok {
					if username == "a" && password == "a" {
						authenticated = true
						SetUser(c, "BasicAuth", username)
						logger.Info().
							Msg("Pass Basic")
					}
				}
			}

			if authenticated {
				logger.Info().
					Msg("Authenticated")
			} else {
				logger.Info().
					Msg("Anonymous")

			}

			return next(c)
		}
	}
}

func returnError(c echo.Context, code int, message string) error {
	var code32 = int32(code)
	errResponse := Error{
		Code:    &code32,
		Message: message,
	}
	return c.JSON(code, errResponse)
}

type ServerImpl struct {
}

func NewServer() ServerInterface {
	return &ServerImpl{}
}

func (s *ServerImpl) GetStatusV1(c echo.Context, params GetStatusV1Params) error {
	// logger := log.Logger
	logger := GetApiLogger(c)
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
	return c.JSON(http.StatusOK, status)
}

func (s *ServerImpl) SetStatusV1(c echo.Context) error {
	// logger := log.Logger
	logger := GetApiLogger(c)
	logger.Info().
		Msg("SetStatusV1")
	var status GetStatusV1
	err := c.Bind(&status)
	if err != nil {
		return returnError(c, http.StatusBadRequest, "could not bind request body")
	}
	logger.Info().
		Interface("value", status).
		Msg("input")
	status.Status = "Ok or Not"
	logger.Info().
		Interface("result", status).
		Msg("output")
	return c.JSON(http.StatusOK, status)
}

type StrictServerImpl struct {
}

func NewStrictServer() StrictServerInterface {
	return &StrictServerImpl{}
}

func (s *StrictServerImpl) GetStatusV1(c context.Context, request GetStatusV1RequestObject) (GetStatusV1ResponseObject, error) {

	logger := log.Logger.With().
		Str("mode", "strict").Logger()

	// в отличии от chi и gin тут в режиме strict доступа к контексту echo не имеется
	ec1 := echoMiddleware.GetEchoContext(c)
	// apiechogen
	ec2 := c.Value(EchoContextKey).(echo.Context)

	logger.Debug().
		Interface("ec1", ec1).
		Interface("ec2", ec2).
		Msg("EchoContext")

	logger.Debug().
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

func (s *StrictServerImpl) SetStatusV1(c context.Context, request SetStatusV1RequestObject) (SetStatusV1ResponseObject, error) {
	logger := log.Logger.With().
		Str("mode", "strict").Logger()

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
