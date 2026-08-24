package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/khaingminhtun/relio-backend/internal/app"
	"github.com/khaingminhtun/relio-backend/internal/config"
	"github.com/khaingminhtun/relio-backend/internal/infrastructure/database"
	redisinfra "github.com/khaingminhtun/relio-backend/internal/infrastructure/redis"
	"github.com/khaingminhtun/relio-backend/internal/shared/email"
	"github.com/khaingminhtun/relio-backend/internal/shared/logger"
	"github.com/khaingminhtun/relio-backend/internal/shared/security"

	"github.com/rs/zerolog/log"
)

func main() {

	// ============================================================
	// Config
	// ============================================================

	cfg := config.Load()

	logger.Init(cfg.Loglevel)

	// ============================================================
	// Database
	// ============================================================

	db, err := database.NewGorm(cfg.DB)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("database initialization failed")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to get sql database")
	}

	// ============================================================
	// Redis
	// ============================================================

	redisClient := redisinfra.NewClient(cfg.Redis)

	ctx := context.Background()

	if err := redisinfra.Ping(ctx, redisClient); err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to connect to redis")
	}

	redisStore := redisinfra.NewStore(redisClient)

	emailQueue := redisinfra.NewEmailQueue(redisClient)

	if err := emailQueue.EnsureConsumerGroup(ctx); err != nil {
		log.Fatal().
			Err(err).
			Msg("email queue initialization failed")
	}

	// ============================================================
	// Email
	// ============================================================

	emailSender := email.NewSendGridSender(
		cfg.SendGrid.APIKey,
		cfg.SendGrid.FromEmail,
		cfg.SendGrid.FromName,
	)

	// ============================================================
	// Email Worker
	// ============================================================

	workerCtx, cancelWorker := context.WithCancel(
		context.Background(),
	)

	emailWorker := email.New(
		emailQueue,
		emailSender,
	)

	go emailWorker.Start(workerCtx)

	log.Info().
		Msg("email worker started")

	// ===============
	// jwt
	// ==============

	jwtManager := security.NewJWTManager(
		cfg.JwtEnv.Secret,
		cfg.JwtEnv.AccessExpiration,
		cfg.JwtEnv.RefreshExpiration,
	)

	// ============================================================
	// Dependency Injection
	// ============================================================

	deps := app.NewDependencies(
		db,
		redisStore,
		emailQueue,

		jwtManager,
	)

	// ============================================================
	// Router
	// ============================================================

	r := app.NewRouter(deps)

	// ============================================================
	// HTTP Server
	// ============================================================

	server := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: r,
	}

	go func() {

		log.Info().
			Str("addr", cfg.ServerPort).
			Msg("server started")

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			log.Fatal().
				Err(err).
				Msg("server failed")
		}
	}()

	// ============================================================
	// Wait for Shutdown Signal
	// ============================================================

	shutdownSignalCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

	<-shutdownSignalCtx.Done()

	log.Info().
		Msg("shutdown signal received")

	// ============================================================
	// Graceful Shutdown
	// ============================================================

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	// Stop accepting new HTTP requests
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error().
			Err(err).
			Msg("http server shutdown failed")
	} else {
		log.Info().
			Msg("http server stopped")
	}

	// Stop background worker
	cancelWorker()

	log.Info().
		Msg("email worker stopped")

	// Close PostgreSQL
	if err := sqlDB.Close(); err != nil {
		log.Error().
			Err(err).
			Msg("database close failed")
	} else {
		log.Info().
			Msg("database connection closed")
	}

	// Close Redis
	if err := redisClient.Close(); err != nil {
		log.Error().
			Err(err).
			Msg("redis close failed")
	} else {
		log.Info().
			Msg("redis connection closed")
	}

	log.Info().
		Msg("application shutdown complete")
}
