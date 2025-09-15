package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Koshsky/polyschedule-backend/internal/config"
	"github.com/Koshsky/polyschedule-backend/pkg/polyschedule"
	"github.com/joho/godotenv"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	if lvl, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL")); err == nil {
		zerolog.SetGlobalLevel(lvl)
	}

	// Загружаем .env (если присутствует), чтобы переменные окружения были доступны как локально, так и в контейнере
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	log.Info().Interface("cfg", cfg).Msg("starting service")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Сервис как модуль (может использоваться из другого репозитория)
	svc := polyschedule.New(cfg)
	httpErrCh, _ := svc.Start(ctx)

	select {
	case <-ctx.Done():
		log.Info().Msg("shutdown signal received")
	case err := <-httpErrCh:
		if err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("http server error")
		}
	}

	// Graceful shutdown
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	svc.Stop(shutdownCtx)

	log.Info().Msg("service stopped")
}
