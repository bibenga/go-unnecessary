package main

import (
	"net/http"
	"os"
	"time"
	"unnecessary/api-chi-gen/server"

	stdLog "log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/unrolled/secure"

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
	logger := log.Logger

	// redirect standart logger
	stdLogWriter := root.With().CallerWithSkipFrameCount(4).Logger()
	stdLog.SetFlags(0)
	stdLog.SetOutput(stdLogWriter)

	// ------------------------------------------------------------------------------------

	secureMiddleware := secure.New(secure.Options{
		SSLRedirect: false,
	})

	// github.com/go-chi/chi
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	// router.Use(middleware.Logger)
	router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  stdLog.New(stdLog.Writer(), "", stdLog.Flags()),
		NoColor: false,
	}))
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(secureMiddleware.Handler)
	router.Use(middleware.NoCache)
	router.Use(middleware.Heartbeat("/ping"))

	router.Mount("/debug", middleware.Profiler())

	fs := http.FileServer(http.Dir("api"))
	router.Mount("/docs2", http.StripPrefix("/docs2/", fs))

	validator, _ := server.NewValidator()

	// api := server.NewServer()
	// router.Mount("/", server.HandlerWithOptions(api, server.ChiServerOptions{
	// 	Middlewares: []server.MiddlewareFunc{
	// 		server.Authenticator,
	// 		validator,
	// 	},
	// }))

	api := server.NewStrictServer()
	strictApi := server.NewStrictHandler(
		api,
		[]server.StrictMiddlewareFunc{
			// 	server.Authenticator,
			// 	validator,
		},
	)
	router.Mount("/", server.HandlerWithOptions(strictApi, server.ChiServerOptions{
		BaseURL: "/api",
		Middlewares: []server.MiddlewareFunc{
			server.Authenticator,
			validator,
		},
	}))

	logger.Info().Msg("Ready")
	http.ListenAndServe(":8000", router)
}
