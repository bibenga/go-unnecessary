// go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest --config=server.cfg.yaml ../api.yaml
// go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.cfg.yaml ../api.yaml
// go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
//
//go:generate oapi-codegen --config=model.cfg.yaml  ../../api/api.yaml
//go:generate oapi-codegen --config=client.cfg.yaml ../../api/api.yaml
//go:generate oapi-codegen --config=server.cfg.yaml ../../api/api.yaml

package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	ginMiddleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
)

func SetRequestId(c *gin.Context, requestId string) {
	c.Set("RequestId", requestId)
}

func GetRequestId(c *gin.Context) string {
	requestId := requestid.Get(c)
	return requestId
}

func SetUser(c *gin.Context, backend, user string) {
	c.Set("AuthenticationBackend", backend)
	c.Set("RequestUser", user)
}

func GetAuthenticationBackend(c *gin.Context) string {
	backendRaw, ok := c.Get("AuthenticationBackend")
	if ok {
		backend, _ := backendRaw.(string)
		return backend
	} else {
		return ""
	}
}

func GetUser(c *gin.Context) string {
	userRaw, ok := c.Get("RequestUser")
	if ok {
		user, _ := userRaw.(string)
		return user
	} else {
		return ""
	}
}

type ApiRequestLoggerHook struct {
	ctx *gin.Context
}

func (h ApiRequestLoggerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	// e.Str("authenticationBackend", GetAuthenticationBackend(h.ctx))
	e.Str("requestId", GetRequestId(h.ctx))
	e.Str("user", GetUser(h.ctx))
}

func SetApitLogger(ctx *gin.Context, parent zerolog.Logger) {
	logger := log.With().
		Logger().Hook(&ApiRequestLoggerHook{ctx: ctx})
	// ctx.Set("ApiLoggerParent", parent)
	ctx.Set("ApiLoggerRequest", logger)
}

func GetApiLogger(ctx *gin.Context) zerolog.Logger {
	loggerRaw, _ := ctx.Get("ApiLoggerRequest")
	return loggerRaw.(zerolog.Logger)
}

func ApiLoggerHandler(root zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		SetApitLogger(c, root)
		c.Next()
	}
}

// gin.HandlerFunc
func Validator() func(c *gin.Context) {
	spec, err := GetSwagger()
	if err != nil {
		panic(fmt.Errorf("loading spec: %w", err))
	}
	validator := ginMiddleware.OapiRequestValidatorWithOptions(spec,
		&ginMiddleware.Options{
			Options: openapi3filter.Options{
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

func AuthenticateAuthenticated(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	c := ginMiddleware.GetGinContext(ctx)
	if c == nil {
		log.Panic().
			Msg("Echo context was not found")
		return ErrUnauthenticated
	}
	// logger := GetApiLogger(c)
	logger := log.Logger

	requestUser := GetUser(c)
	if requestUser == "" {
		logger.Info().Msgf("validator: rejected")
		return errors.New("ErrUnauthenticated")
	} else {
		logger.Info().Msgf("validator: pass")
		return nil
	}
}

func Authenticator(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authenticated := false

		if !authenticated {
			authHdr := c.Request.Header.Get("X-API-Key")
			logger.Info("Try ApiKey",
				zap.String("Header", "X-API-Key"),
				zap.String("Value", authHdr))
			if authHdr == "a" {
				authenticated = true
				SetUser(c, "ApiKeyAuth", authHdr)
				logger.Info("Pass ApiKey")
			}
		}

		if !authenticated {
			username, password, ok := c.Request.BasicAuth()
			logger.Info("Try Basic",
				zap.Bool("ok", ok),
				zap.String("username", username),
				zap.String("password", password))
			if ok {
				if username == "a" && password == "a" {
					authenticated = true
					SetUser(c, "BasicAuth", username)
					logger.Info("Pass Basic")
				}
			}
		}

		if authenticated {
			logger.Info("Authenticated")
		}

		c.Next()
	}
}

func returnError(ctx *gin.Context, code int, message string) {
	var code32 = int32(code)
	errResponse := Error{
		Code:    &code32,
		Message: message,
	}
	ctx.JSON(code, errResponse)
}

type ServerImpl struct {
	logger *zap.Logger
}

func NewServer(logger *zap.Logger) ServerInterface {
	return &ServerImpl{logger: logger}
}

func (s *ServerImpl) GetStatusV1(c *gin.Context, params GetStatusV1Params) {
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
	c.JSON(http.StatusOK, status)
}

func (s *ServerImpl) SetStatusV1(c *gin.Context) {
	// logger := log.Logger
	logger := GetApiLogger(c)
	logger.Info().
		Msg("SetStatusV1")
	var status GetStatusV1
	err := c.Bind(&status)
	if err != nil {
		returnError(c, http.StatusBadRequest, "could not bind request body")
		return
	}
	logger.Info().
		Interface("value", status).
		Msg("input")
	status.Status = "Ok or Not"
	logger.Info().
		Interface("result", status).
		Msg("output")
	c.JSON(http.StatusOK, status)
}

type StrictServerImpl struct {
	logger *zap.Logger
}

func NewStrictServer(logger *zap.Logger) StrictServerInterface {
	return &StrictServerImpl{
		logger: logger,
	}
}

func (s *StrictServerImpl) GetStatusV1(c context.Context, request GetStatusV1RequestObject) (GetStatusV1ResponseObject, error) {
	ctx := c.(*gin.Context)
	logger := s.logger.With(
		zap.String("requestId", GetRequestId(ctx)),
		zap.String("user", GetUser(ctx)),
	)

	logger.Info("GetStatusV1", zap.Reflect("input", request))
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
	logger.Info("result", zap.Reflect("output", result))
	return result, nil
}

func (s *StrictServerImpl) SetStatusV1(c context.Context, request SetStatusV1RequestObject) (SetStatusV1ResponseObject, error) {
	ctx := c.(*gin.Context)
	logger := s.logger.With(
		zap.String("requestId", GetRequestId(ctx)),
		zap.String("user", GetUser(ctx)),
	)

	logger.Info("SetStatusV1", zap.Reflect("input", request))
	result := SetStatusV1200JSONResponse{
		Status: "Ok or Not",
	}
	logger.Info("result", zap.Reflect("output", result))
	return result, nil
}
