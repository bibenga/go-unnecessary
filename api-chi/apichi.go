package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ErrUnauthenticated = errors.New("ErrUnauthenticated")
	ErrApiKey          = errors.New("ErrApiKey")
	ErrBasic           = errors.New("ErrBasic")
)

type ApiError struct {
	Code      int     `json:"code" validate:"required" example:"400"`
	Message   string  `json:"message" validate:"required" example:"Invalid Request"`
	RequestID *string `json:"requestId" example:"slfjlskjfglksjfg"`
}

type Status struct {
	Status string `json:"status" validate:"required" example:"OK"`
}

type Handler struct {
	logger *zap.Logger
}

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("process GetStatus")
	isFull := r.URL.Query().Get("isFull") == "true"
	q := r.URL.Query().Get("q")

	pageRaw := r.URL.Query().Get("page")
	page := 1
	if pageRaw != "" {
		page1, err := strconv.Atoi(pageRaw)
		if err != nil {
			h.logger.Warn("Can't read page number",
				zap.Error(err))
		} else {
			page = page1
		}
	}
	h.logger.Info("input",
		zap.Bool("isFull", isFull),
		zap.String("q", q),
		zap.Int("page", page))

	status := &Status{
		Status: "OK",
	}

	h.logger.Info("output",
		zap.Any("value", status))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, status)
}

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

	// github.com/go-chi/chi
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	// router.Use(middleware.RequestLogger(
	// 	&middleware.DefaultLogFormatter{
	// 		Logger:  log.New(logger, "", 0),
	// 		NoColor: true,
	// 	}))
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(middleware.NoCache)
	router.Use(middleware.Heartbeat("/ping"))

	// router.Use(hlog.NewHandler(log.Logger))
	// router.Use(hlog.RequestIDHandler("requestID", "X-Request-ID"))
	// router.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
	// 	logger.Info(r.RequestURI,
	// 		zap.String("method", r.Method),
	// 		zap.Int("status", status),
	// 		zap.Int("responseSize", size),
	// 		zap.Duration("duration", duration))
	// }))

	router.Mount("/debug", middleware.Profiler())

	fs := http.FileServer(http.Dir("api"))
	router.Mount("/docs2", http.StripPrefix("/docs2/", fs))

	h := Handler{logger: logger}
	router.Get("/api/v1/status", h.GetStatus)
	// root.Post("/api/v1/status", SetStatus)

	logger.Info("Ready on port 8000: http://127.0.0.1:8000/docs2/swagger.html")
	if err := http.ListenAndServe(":8000", router); err != nil {
		panic(err)
	}
}
