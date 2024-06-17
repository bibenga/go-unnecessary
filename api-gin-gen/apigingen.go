package main

import (
	"log"
	"net/http"
	"unnecessary/api-gin-gen/server"

	ginZero "github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-contrib/secure"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	ginZap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.ConsoleSeparator = " "
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	zap.RedirectStdLog(logger)

	// ------------------------------------------------------------------------------------
	gin.SetMode(gin.DebugMode)

	r := gin.New()

	r.Use(requestid.New())

	r.Use(ginZero.SetLogger())
	r.Use(ginZap.GinzapWithConfig(logger.Named("access"), &ginZap.Config{
		TimeFormat: "",
		UTC:        true,
		Context: func(c *gin.Context) []zapcore.Field {
			return []zapcore.Field{
				zap.String("requestID", server.GetRequestId(c)),
				zap.String("requestUser", server.GetUser(c)),
				zap.String("authenticationBackend", server.GetAuthenticationBackend(c)),
			}
		},
	}))

	secureConfig := secure.DefaultConfig()
	secureConfig.SSLRedirect = false
	secureConfig.IsDevelopment = true
	r.Use(secure.New(secureConfig))
	r.Use(gin.Recovery())

	store := cookie.NewStore([]byte("supersecret"))
	r.Use(sessions.Sessions("gin-session", store))

	// r.Use(server.ApiLoggerHandler(log.Logger))
	r.Use(server.Authenticator(logger.Named("auth")))
	// r.Use(server.Validator())

	// h2 := server.NewServer()
	// server.RegisterHandlersWithOptions(r, h2, server.GinServerOptions{
	// 	Middlewares: []server.MiddlewareFunc{
	// 		server.Validator(),
	// 	},
	// })

	ssrv2 := server.NewStrictServer(logger.Named("api"))
	h2 := server.NewStrictHandler(ssrv2, []server.StrictMiddlewareFunc{})
	server.RegisterHandlersWithOptions(r, h2, server.GinServerOptions{
		BaseURL: "/api",
		Middlewares: []server.MiddlewareFunc{
			server.Validator(),
		},
	})

	r.Static("/docs2", "api")

	log.Printf("Ready or not")
	logger.Info("Ready")

	s := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}
	if err := s.ListenAndServe(); err != nil {
		logger.Panic("Failed", zap.Error(err))
	} else {
		logger.Panic("Terminated")
	}
}
