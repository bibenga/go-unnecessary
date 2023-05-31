package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	stdLog "log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

// =======================

type UnnecessaryValidator struct {
	validator *validator.Validate
}

func (cv *UnnecessaryValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

type ApiError struct {
	Code      int     `json:"code" validate:"required" example:"400"`
	Message   string  `json:"message" validate:"required" example:"Invalid Request"`
	RequestID *string `json:"requestId" example:"slfjlskjfglksjfg"`
}

var (
	ErrUnauthenticated = errors.New("ErrUnauthenticated")
	ErrApiKey          = errors.New("ErrApiKey")
	ErrBasic           = errors.New("ErrBasic")
)

func Authenticate(logger zerolog.Logger, ctx echo.Context) error {
	logger.Info().
		Msg("Authenticate")

	// ApiKey
	logger.Info().
		Msg("Try ApiKey")
	authHdrKey := http.CanonicalHeaderKey("X-API-Key")
	authHdr := ctx.Request().Header.Get(authHdrKey)
	logger.Info().
		Str("Header", "X-API-Key").
		Str("Value", authHdr).
		Msg("Try authenenticate user by apiKey")
	if authHdr != "" {
		if authHdr == "a" {
			SetRequestUser(ctx, authHdr)
			logger.Info().
				Msg("Pass ApiKey")
			return nil
		}
		logger.Info().
			Msg("Reject ApiKey")
		return ErrApiKey
	}

	// Basic
	logger.Info().
		Msg("Try Basic")
	username, password, ok := ctx.Request().BasicAuth()
	logger.Info().
		Bool("ok", ok).
		Str("username", username).
		Str("password", "******").
		Msg("Try authenenticate user by username")
	if ok {
		if username == "a" && password == "a" {
			SetRequestUser(ctx, username)
			logger.Info().
				// Caller().
				Msg("Pass Basic")
			return nil
		}
		logger.Info().
			Msg("Reject basic")
		return ErrBasic
	}

	logger.Info().
		Msg("Unauthenticated")
	return ErrUnauthenticated
}

func SaveRequestID(c echo.Context, requestId string) {
	c.Set("RequestID", requestId)
}

func GetRequestID(c echo.Context) string {
	requestID, ok := c.Get("RequestID").(string)
	if ok {
		return requestID
	} else {
		return ""
	}
}

func SetRequestUser(c echo.Context, requestUser string) {
	c.Set("RequestUser", requestUser)
}

func GetRequestUser(c echo.Context) string {
	requestUser, ok := c.Get("RequestUser").(string)
	if ok {
		return requestUser
	} else {
		return ""
	}
}

type ApiRequestLoggerHook struct {
	ctx echo.Context
}

func (h ApiRequestLoggerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	// e.Str("requestID", GetRequestID(h.ctx))
	e.Str("requestUser", GetRequestUser(h.ctx))
}

func GetApiLogger(ctx echo.Context) zerolog.Logger {
	logger, ok := ctx.Get("ApiLogger").(zerolog.Logger)
	if ok {
		return logger
	}

	logger = log.Logger.With().
		Str("requestID", GetRequestID(ctx)).
		// Str("requestUser", GetRequestUser(ctx)).
		Logger()
	// logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
	// 	return c.
	// 		Str("requestID", GetRequestID(ctx)).
	// 		Str("requestUser", GetRequestUser(ctx))
	// })
	// Только хуки работают на лету, UpdateContext работает в момент вызова
	logger = logger.Hook(&ApiRequestLoggerHook{ctx: ctx})

	ctx.Set("ApiLogger", logger)
	return logger
}

type Status struct {
	Status string `json:"status" validate:"required" example:"OK"`
}

// Status godoc
// @Summary      Get status
// @Description  Get service status
// @Tags         status
// @Produce      json
// @Param        isFull	query     boolean  false "IsFull status or not"
// @Param        q    	query     string  false  "search by q"
// @Param        page   query     string  false  "page"
// @Success      200  {object}  Status
// @Failure      400  {object}  ApiError
// @Failure      401  {object}  ApiError
// @Failure      403  {object}  ApiError
// @Failure      500  {object}  ApiError
// @Router       /api/v1/status [get]
func GetStatus(ctx echo.Context) (err error) {
	logger := GetApiLogger(ctx)
	// logger := log.Logger

	isFull := ctx.QueryParam("isFull") == "true"
	q := ctx.QueryParam("q")

	pageRaw := ctx.QueryParam("page")
	page := 1
	if pageRaw != "" {
		page, err = strconv.Atoi(pageRaw)
		if err != nil {
			logger.Warn().
				// Str("requestId", GetRequestID(ctx)).
				Err(err).
				Msg("Can't read page number")
		}
	}

	logger.Info().
		// Str("requestId", GetRequestID(ctx)).
		Bool("isFull", isFull).
		Str("q", q).
		Int("page", page).
		Msg("GetStatus")

	status := &Status{
		Status: "OK",
	}

	logger.Info().
		// Str("requestId", GetRequestID(ctx)).
		Interface("result", status).
		Msg("return")
	return ctx.JSON(http.StatusOK, status)
}

type SetStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

type SetStatusResponse struct {
	SetStatusRequest
	IsSuccess bool `json:"isSuccess"`
}

// Status godoc
// @Security 	BasicAuth || ApiKeyAuth
// @Summary     Set status
// @Description Set service status
// @Tags        status
// @Accept      json
// @Produce     json
// @Param 		request body SetStatusRequest true "query params"
// @Success     200 {object} SetStatusResponse
// @Router      /api/v1/status [post]
func SetStatus(ctx echo.Context) (err error) {
	logger := GetApiLogger(ctx)
	// logger := log.Logger

	logger.Info().
		Msg("process SetStatus")
	if err = Authenticate(logger, ctx); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	var requestStatus SetStatusRequest
	if err = ctx.Bind(&requestStatus); err != nil {
		logger.Warn().
			Err(err).
			Msg("request is bad")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	logger.Info().
		Msgf("parsed request: %+v", requestStatus)
	if err = ctx.Validate(requestStatus); err != nil {
		logger.Warn().
			Err(err).
			Msg("request is invalid")
		return err
	}
	logger.Info().
		Msgf("validated request: %+v", requestStatus)

	var responseStatus SetStatusResponse
	responseStatus.Status = fmt.Sprintf("-> %s <-", requestStatus.Status)

	logger.Info().
		Msgf("response: %+v", responseStatus)
	return ctx.JSON(http.StatusOK, responseStatus)
}

// @title Swagger Example API
// @version 1.0
// @description This is an unnecessary server.
// @contact.name API Support
// @host 127.0.0.1:8000
// //@BasePath /api

// @securityDefinitions.basic  BasicAuth

// @securityDefinitions.apikey  ApiKeyAuth
// @in  header
// @name  X-API-Key
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zoutput := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339Nano,
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.CallerFieldName,
			zerolog.MessageFieldName,
		},
	}
	root := log.Output(zoutput)
	log.Logger = root.With().Caller().Logger()

	// redirect standart logger
	stdLogWriter := root.With().CallerWithSkipFrameCount(4).Logger()
	stdLog.SetFlags(0)
	stdLog.SetOutput(stdLogWriter)

	stdLog.Print("Kaka1")
	log.Info().Msg("Unnecessary1")

	// ------------------------------------------------------------------------------------
	e := echo.New()
	e.Debug = true
	e.Logger.SetOutput(root.With().CallerWithSkipFrameCount(4).Logger())

	e.Validator = &UnnecessaryValidator{validator: validator.New()}

	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		RequestIDHandler: SaveRequestID,
	}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:       true,
		LogURI:          true,
		LogStatus:       true,
		LogResponseSize: true,
		LogRequestID:    true,
		LogLatency:      true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger := GetApiLogger(c)
			logger.Info().
				// Str("requestID", v.RequestID).
				// Str("requestUser", GetRequestUser(c)).
				// Str("method", v.Method).
				Int("status", v.Status).
				Dur("latency", v.Latency).
				Int64("responseSize", v.ResponseSize).
				// Str("uri", v.URI).
				Msgf("%s %s", v.Method, v.URI)
			return nil
		},
	}))

	e.GET("/docs2", func(ctx echo.Context) (err error) {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/docs2/")
	})
	e.GET("/docs2/", func(ctx echo.Context) (err error) {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/docs2/swagger.html")
	})
	e.GET("/docs2/index.html", func(ctx echo.Context) (err error) {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/docs2/swagger.html")
	})
	e.Static("/docs2", "api")
	e.GET("/api/v1/status", GetStatus)
	e.POST("/api/v1/status", SetStatus)
	if err := e.Start(":8000"); err != nil {
		log.Panic().
			Err(err).
			Msg("Terminated")
	} else {
		log.Info().
			Msg("Terminated")
	}
}
