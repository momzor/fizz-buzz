// Package main is the composition root of the fizzbuzz API: it wires
// configuration, adapters and use cases together and starts the HTTP server.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	httpadapter "fizzbuzz/internal/adapter/http"
	mongorepo "fizzbuzz/internal/adapter/repository/mongo"
	"fizzbuzz/internal/config"
	"fizzbuzz/internal/telemetry"
	"fizzbuzz/internal/usecase"

	_ "fizzbuzz/docs" // swagger generated docs
)

// @title           FizzBuzz API
// @version         1.0
// @description     FizzBuzz REST API with usage statistics.
// @BasePath        /
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := run(logger); err != nil {
		logger.Error("service exited with error", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	_ = godotenv.Load()

	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdownTracing, err := telemetry.Setup(ctx, cfg.ServiceName, cfg.ServiceVersion, cfg.OTLPEndpoint, cfg.TracingEnabled)
	if err != nil {
		return err
	}
	defer func() { _ = shutdownTracing(context.Background()) }()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return err
	}
	defer func() { _ = mongoClient.Disconnect(context.Background()) }()

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := mongoClient.Ping(pingCtx, nil); err != nil {
		return err
	}

	statsRepo := mongorepo.NewStatsRepository(mongoClient.Database(cfg.MongoDatabase))
	if err := statsRepo.EnsureIndexes(ctx); err != nil {
		return err
	}

	fizzBuzzUseCase := usecase.NewFizzBuzz(statsRepo, logger)
	statsUseCase := usecase.NewStats(statsRepo)

	fizzBuzzHandler := httpadapter.NewFizzBuzzHandler(fizzBuzzUseCase)
	statsHandler := httpadapter.NewStatsHandler(statsUseCase)

	router := httpadapter.NewRouter(fizzBuzzHandler, statsHandler)

	server := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("starting HTTP server", "addr", cfg.HTTPAddr, "base_url", cfg.PublicBaseURL, "environment", cfg.Environment)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	return server.Shutdown(shutdownCtx)
}
