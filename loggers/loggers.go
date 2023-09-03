package main

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func playStdLog() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")
	log.Printf("------------")
	log.Printf("playStdLog: %v", 1)

	l := slog.New(slog.Default().Handler()).With(slog.String("func", "playStdLog"))
	l.Info("message from slog",
		slog.Int("count", 3),
		slog.Group("request", slog.String("method", "GET"), slog.Int("status", 400)),
	)
}

func playZerolog() {
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
	root := zlog.Output(zoutput)
	zlog.Logger = root.With().Caller().Logger()
	logger := zlog.Logger

	// redirect standart logger
	stdLogWriter := root.With().CallerWithSkipFrameCount(4).Logger()
	log.SetFlags(0)
	log.SetOutput(stdLogWriter)

	logger.Info().Msg("------------")
	logger.Info().
		Float64("float", 63).
		Str("aaaa", "bbbb").
		Msg("playZerolog 1")
	logger.Info().
		Float64("float", 63).
		Str("aaaa", "bbbb").
		Msgf("playZerolog %v", 2)
	log.Printf("playZerolog %v", 3)
}

func playZap() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.ConsoleSeparator = " "
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	zap.RedirectStdLog(logger.Named("std"))

	logger.Info("------------")
	logger.Info("playZap 1", zap.Float64("float", 62.1), zap.String("aaaa", "cccc"))
	logger.Sugar().Infow("playZap 23", "float", 62.1, "aaaa", "cccc")
	logger.Sugar().Infof("playZap %v", 21, zap.Float64("float", 62.1), zap.String("aaaa", "cccc"))
	logger.Sugar().Infof("playZap %v", 22, 62.1, "cccc")
	log.Printf("playZap %v", 3)
}

func main() {
	playStdLog()
	playZerolog()
	playZap()
}
