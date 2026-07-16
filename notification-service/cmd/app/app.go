package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anomalyco/hookah-store/notification-service/internal/config"
	"github.com/anomalyco/hookah-store/notification-service/internal/repository/postgres/notification"
	"github.com/anomalyco/hookah-store/notification-service/internal/services/admin"
	"github.com/anomalyco/hookah-store/notification-service/internal/transport/http"
	adminHandlers "github.com/anomalyco/hookah-store/notification-service/internal/transport/http/handlers/admin"
	"github.com/anomalyco/hookah-store/notification-service/pkg/database"
)

const (
	serviceName = "notification-service"
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
		slog.Error("failed to initialize config", slog.String("err", err.Error()))
		os.Exit(1)
	}

	slog.Info("config loaded", slog.String("env", cfg.Env))

	db, err := database.NewDB(&cfg.DataBase)
	if err != nil {
		slog.Error("failed to connect to database", slog.String("err", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("database connected", slog.String("db_name", cfg.DataBase.DBName))

	notifRepo := notification.New(db)
	adminService := admin.New(notifRepo)
	adminHandler := adminHandlers.New(adminService)

	httpServer := http.New(&cfg.HTTPServer, adminHandler)
	go func() {
		if err := httpServer.Run(); err != nil {
			slog.Error("failed to start http server", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	slog.Info("received shutdown signal", slog.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("err", err.Error()))
	}

	slog.Info("service stopped", slog.String("service", serviceName))
}
