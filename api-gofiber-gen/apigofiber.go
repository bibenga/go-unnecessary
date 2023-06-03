package main

import (
	"net/http"
	"os"
	"time"
	"unnecessary/api-gofiber-gen/server"

	stdLog "log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
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
	log.Info().Msg("Unnecessary1")

	// ------------------------------------------------------------------------------------
	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(requestid.New())
	app.Use(helmet.New())
	app.Use(etag.New())

	api := app.Group("/api")
	api.Use(server.Authenticator())
	api.Use(server.Validator())

	ssrv2 := server.NewStrictServer()
	h2 := server.NewStrictHandler(ssrv2, nil)
	server.RegisterHandlersWithOptions(api, h2, server.FiberServerOptions{
		// BaseURL: "/api",
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/docs2/swagger.html", http.StatusTemporaryRedirect)
	})
	app.Get("/docs2", func(c *fiber.Ctx) error {
		return c.Redirect("/docs2/swagger.html", http.StatusTemporaryRedirect)
	})
	app.Get("/docs2/", func(c *fiber.Ctx) error {
		return c.Redirect("/docs2/swagger.html", http.StatusTemporaryRedirect)
	})
	app.Get("/docs2/index.html", func(c *fiber.Ctx) error {
		return c.Redirect("/docs2/swagger.html", http.StatusTemporaryRedirect)
	})
	app.Static("/docs2", "api")

	if err := app.Listen(":8000"); err != nil {
		log.Panic().
			Err(err).
			Msg("Terminated")
	} else {
		log.Info().
			Msg("Terminated")
	}
}
