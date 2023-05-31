package main

import (
	"net/http"
	"os"
	"time"
	"unnecessary/api-echo-gen/server"

	stdLog "log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

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
	log.Info().Msg("Unnecessary2")

	// ------------------------------------------------------------------------------------
	// https://github.com/deepmap/oapi-codegen
	e := echo.New()
	e.Debug = true
	e.Logger.SetOutput(root.With().CallerWithSkipFrameCount(4).Logger())

	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		RequestIDHandler: func(ctx echo.Context, s string) {
			server.SetRequestId(ctx, s)
		},
	}))
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:       true,
		LogURI:          true,
		LogStatus:       true,
		LogResponseSize: true,
		LogRequestID:    true,
		LogLatency:      true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			// logger := log.Logger
			logger := server.GetApiLogger(c)
			logger.Info().
				Str("requestId", v.RequestID).
				// Str("method", v.Method).
				Int("status", v.Status).
				Dur("latency", v.Latency).
				Int64("responseSize", v.ResponseSize).
				// Str("uri", v.URI).
				Msgf("%s %s", v.Method, v.URI)
			return nil
		},
	}))

	e.Use(server.ApiLoggerHandler(log.Logger))
	e.Use(server.StoreEchoContextInRequestContext())
	e.Use(server.Authenticator())
	e.Use(server.Validator())

	// svr := server.NewServer()
	// server.RegisterHandlers(e, svr)

	ssrv2 := server.NewStrictServer()
	h2 := server.NewStrictHandler(ssrv2, nil)
	server.RegisterHandlers(e, h2)

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

	if err := e.Start(":8000"); err != nil {
		log.Panic().
			Err(err).
			Msg("Terminated")
	} else {
		log.Info().
			Msg("Terminated")
	}
}
