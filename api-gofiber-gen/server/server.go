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

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	fiberMiddleware "github.com/oapi-codegen/fiber-middleware"
	"github.com/rs/zerolog/log"
)

// func SetUser(c *fiber.Ctx, backend, user string) {
// 	c.Set("AuthenticationBackend", backend)
// 	c.Set("RequestUser", user)
// }

// func GetAuthenticationBackend(c *fiber.Ctx) string {
// 	backend, ok := c.UserContext().Get("AuthenticationBackend")
// 	if ok {
// 		return backend
// 	} else {
// 		return ""
// 	}
// }

// func GetUser(c echo.Context) string {
// 	user, ok := c.Get("RequestUser").(string)
// 	if ok {
// 		return user
// 	} else {
// 		return ""
// 	}
// }

// type ApiRequestLoggerHook struct {
// 	ctx echo.Context
// }

// func (h ApiRequestLoggerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
// 	// e.Str("authenticationBackend", GetAuthenticationBackend(h.ctx))
// 	e.Str("requestId", GetRequestId(h.ctx))
// 	e.Str("user", GetUser(h.ctx))
// }

// func SetApitLogger(ctx echo.Context, parent zerolog.Logger) {
// 	// logger := log.With().
// 	// 	Str("requestId", GetRequestId(ctx)).
// 	// 	Str("user", user).
// 	// 	Logger().Hook()
// 	// logger := log.With().
// 	// 	Logger().Hook(&ApiLoggerHook{ctx: ctx})
// 	logger := log.With().
// 		Logger().Hook(&ApiRequestLoggerHook{ctx: ctx})
// 	// ctx.Set("ApiLoggerParent", parent)
// 	ctx.Set("ApiLoggerRequest", logger)
// }

// func GetApiLogger(ctx echo.Context) zerolog.Logger {
// 	return ctx.Get("ApiLoggerRequest").(zerolog.Logger)
// }

// func ApiLoggerHandler(root zerolog.Logger) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(ctx echo.Context) error {
// 			SetApitLogger(ctx, root)
// 			return next(ctx)
// 		}
// 	}
// }

// type contextKey string

// const EchoContextKey = contextKey("unnecessary/EchoContext")

// func StoreEchoContextInRequestContext() echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			r := c.Request()
// 			ctx := r.Context()
// 			ctx = context.WithValue(ctx, EchoContextKey, c)
// 			c.SetRequest(r.WithContext(ctx))
// 			return next(c)
// 		}
// 	}
// }

func Validator() fiber.Handler {
	spec, err := GetSwagger()
	if err != nil {
		panic(fmt.Errorf("loading spec: %w", err))
	}
	validator := fiberMiddleware.OapiRequestValidatorWithOptions(spec,
		&fiberMiddleware.Options{
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
	// c := echoMiddleware.GetEchoContext(ctx)
	// if c == nil {
	// 	log.Panic().
	// 		Msg("Echo context was not found")
	// 	return ErrUnauthenticated
	// }
	// logger := GetApiLogger(c)
	// // logger := log.Logger

	// requestUser := GetUser(c)
	// if requestUser == "" {
	// 	logger.Info().Msgf("validator: rejected")
	// 	return errors.New("ErrUnauthenticated")
	// } else {
	// 	logger.Info().Msgf("validator: pass")
	// 	return nil
	// }
	// logger.Info().Msgf("validator: pass")
	return nil
}

func Authenticator() fiber.Handler {
	logger := log.Logger
	return func(c *fiber.Ctx) error {
		authenticated := false
		if !authenticated {
			authHdr := c.Get("X-API-Key")
			logger.Info().
				Str("Header", "X-API-Key").
				Str("Value", authHdr).
				Msg("Try ApiKey")
			if authHdr == "a" {
				authenticated = true
				c.Locals("ApiKeyAuth", authHdr)
				logger.Info().
					Msg("Pass ApiKey")
			}
		}

		// if !authenticated {
		// 	username, password, ok := c.Request().BasicAuth()
		// 	logger.Info().
		// 		Bool("ok", ok).
		// 		Str("username", username).
		// 		Str("password", password).
		// 		Msg("Try Basic")
		// 	if ok {
		// 		if username == "a" && password == "a" {
		// 			authenticated = true
		// 			SetUser(c, "BasicAuth", username)
		// 			logger.Info().
		// 				Msg("Pass Basic")
		// 		}
		// 	}
		// }

		if authenticated {
			logger.Info().
				Msg("Authenticated")
		} else {
			logger.Info().
				Msg("Anonymous")
		}
		return c.Next()
	}
}

type ServerImpl struct {
}

func NewServer() ServerInterface {
	return &ServerImpl{}
}

func (s *ServerImpl) GetStatusV1(c *fiber.Ctx, params GetStatusV1Params) error {
	logger := log.Logger
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
	return c.JSON(status)
}

func (s *ServerImpl) SetStatusV1(c *fiber.Ctx) error {
	logger := log.Logger
	// logger := GetApiLogger(c)
	logger.Info().
		Msg("SetStatusV1")
	var status GetStatusV1
	if err := c.BodyParser(&status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	logger.Info().
		Interface("value", status).
		Msg("input")
	status.Status = "Ok or Not"
	logger.Info().
		Interface("result", status).
		Msg("output")
	return c.JSON(status)
}

type StrictServerImpl struct {
	store *session.Store
}

func NewStrictServer(store *session.Store) StrictServerInterface {
	return &StrictServerImpl{
		store: store,
	}
}

func (s *StrictServerImpl) GetStatusV1(c context.Context, request GetStatusV1RequestObject) (GetStatusV1ResponseObject, error) {
	logger := log.Logger.With().
		Str("mode", "strict").Logger()

	ctx := c.Value("uFiberContext").(*fiber.Ctx)
	sess, err := s.store.Get(ctx)
	if err != nil {
		return nil, err
	}
	sess.Set("Olala", "1")
	if err := sess.Save(); err != nil {
		return nil, err
	}

	// в отличии от chi и gin тут в режиме strict доступа к контексту echo не имеется
	// ec1 := echoMiddleware.GetEchoContext(c)
	// apiechogen
	// ec2 := c.Value(EchoContextKey).(echo.Context)

	// logger.Debug().
	// 	Interface("ec1", ec1).
	// 	Interface("ec2", ec2).
	// 	Msg("EchoContext")

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
