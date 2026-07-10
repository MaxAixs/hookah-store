package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anomalyco/hookah-store/user-service/internal/config"
	"github.com/anomalyco/hookah-store/user-service/internal/repository/postgres"
	"github.com/anomalyco/hookah-store/user-service/internal/services"
	"github.com/anomalyco/hookah-store/user-service/internal/transport/http"
	"github.com/anomalyco/hookah-store/user-service/internal/transport/http/handlers/admin"
	"github.com/anomalyco/hookah-store/user-service/internal/transport/http/handlers/auth"
	"github.com/anomalyco/hookah-store/user-service/pkg/database"
	jwtpkg "github.com/anomalyco/hookah-store/user-service/pkg/jwt"
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
		slog.Error("failed to connect to database", err)
		os.Exit(1)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepo(db)

	jwtCfg := jwtpkg.New(cfg.JWT.Secret, cfg.JWT.TTL)
	authService := services.NewAuth(userRepo, jwtCfg)
	userService := services.NewAdmin(userRepo)

	adminHandlers := admin.New(userService)
	authHandlers := auth.New(authService)

	httpServer := http.New(&cfg.HTTPServer, jwtCfg, authHandlers, adminHandlers)
	go func() {
		if err := httpServer.Run(); err != nil {
			slog.Error("failed to start http server", "err", err)
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
		slog.Error("failed to shutdown server", "err", err)
		os.Exit(1)
	}

	slog.Info("service stopped", slog.String("service", serviceName))
}
