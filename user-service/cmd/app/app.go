package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anomalyco/hookah-store/user-service/internal/config"
	"github.com/anomalyco/hookah-store/user-service/internal/relay"
	postgresoutbox "github.com/anomalyco/hookah-store/user-service/internal/repository/postgres/outbox"
	postgresuser "github.com/anomalyco/hookah-store/user-service/internal/repository/postgres/user"
	authservice "github.com/anomalyco/hookah-store/user-service/internal/services/auth"
	userservice "github.com/anomalyco/hookah-store/user-service/internal/services/user"
	"github.com/anomalyco/hookah-store/user-service/internal/transport/http"
	"github.com/anomalyco/hookah-store/user-service/internal/transport/http/handlers/admin"
	authhandler "github.com/anomalyco/hookah-store/user-service/internal/transport/http/handlers/auth"
	"github.com/anomalyco/hookah-store/user-service/pkg/database"
	jwtpkg "github.com/anomalyco/hookah-store/user-service/pkg/jwt"
	kafkapkg "github.com/anomalyco/hookah-store/user-service/pkg/kafka"
)

const (
	serviceName = "user-service"
)

func Start() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info(
		"service starting",
		slog.String("service", serviceName),
	)

	cfg, err := config.New()
	if err != nil {
		slog.Error("failed to initialize config", "err", err)
		os.Exit(1)
	}

	slog.Info("config loaded", slog.String("env", cfg.Env))

	db, err := database.NewDB(&cfg.DataBase)
	if err != nil {
		slog.Error("failed to connect to database", slog.String("err", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	userRepo := postgresuser.New(db)
	outboxRepo := postgresoutbox.New(db)

	jwtCfg := jwtpkg.New(cfg.JWT.Secret, cfg.JWT.TTL)

	publisher := kafkapkg.NewPublisher(cfg.Kafka)
	relaySrv := relay.NewOutboxRelay(outboxRepo, publisher)

	authService := authservice.New(db, userRepo, jwtCfg)
	userService := userservice.New(userRepo)

	adminHandlers := admin.New(userService)
	authHandlers := authhandler.New(authService)

	httpServer := http.New(&cfg.HTTPServer, jwtCfg, authHandlers, adminHandlers)
	go func() {
		if err := httpServer.Run(); err != nil {
			slog.Error("failed to start http server", "err", err)
			os.Exit(1)
		}
	}()

	relayCtx, relayCancel := context.WithCancel(context.Background())
	go func() {
		if err := relaySrv.Run(relayCtx); err != nil && err != context.Canceled {
			slog.Error("outbox relay failed", slog.String("err", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	slog.Info("received shutdown signal", slog.String("signal", sig.String()))

	relayCancel()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("err", err.Error()))
	}

	if err := publisher.Close(); err != nil {
		slog.Error("failed to close kafka publisher", slog.String("err", err.Error()))
	}

	slog.Info("service stopped", slog.String("service", serviceName))
}
