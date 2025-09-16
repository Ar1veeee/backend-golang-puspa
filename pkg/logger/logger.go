package logger

import (
	"backend-golang/internal/infrastructure/config"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	appEnv := config.GetEnv("APP_ENV", "development")

	var logger zerolog.Logger
	if appEnv == "development" {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	} else {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	log.Logger = logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Info().Msgf("Logger initialized for environment: %s", appEnv)
}
