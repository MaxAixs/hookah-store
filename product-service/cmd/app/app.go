package app

import (
	"context"
	"errors"
	"log/slog"
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anomalyco/hookah-store/product-service/internal/config"
	postgrescategory "github.com/anomalyco/hookah-store/product-service/internal/repository/postgres/category"
	postgresproduct "github.com/anomalyco/hookah-store/product-service/internal/repository/postgres/product"
	categoryservice "github.com/anomalyco/hookah-store/product-service/internal/services/category"
	inventoryservice "github.com/anomalyco/hookah-store/product-service/internal/services/inventory"
	productservice "github.com/anomalyco/hookah-store/product-service/internal/services/product"
	"github.com/anomalyco/hookah-store/product-service/internal/transport/http"
	adminhandlers "github.com/anomalyco/hookah-store/product-service/internal/transport/http/handlers/admin"
	userhandlers "github.com/anomalyco/hookah-store/product-service/internal/transport/http/handlers/user"
	"github.com/anomalyco/hookah-store/product-service/pkg/database"
	jwtpkg "github.com/anomalyco/hookah-store/user-service/pkg/jwt"
)

const (
	serviceName = "product-service"
)

func Start() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("service starting", slog.String("service", serviceName))

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

	categoryRepo := postgrescategory.New(db)
	productRepo := postgresproduct.New(db)

	jwtCfg := jwtpkg.New(cfg.JWT.Secret, cfg.JWT.TTL)

	categoryService := categoryservice.New(categoryRepo)
	productService := productservice.New(db, productRepo, productRepo)
	inventoryService := inventoryservice.New(productRepo)

	adminHandler := adminhandlers.New(categoryService, productService, inventoryService)
	userHandler := userhandlers.New(categoryService, productService)

	httpServer := http.New(&cfg.HTTPServer, jwtCfg,
		[]http.Handler{adminHandler},
		[]http.PublicHandler{userHandler},
	)
	go func() {
		if err := httpServer.Run(); err != nil && !errors.Is(err, stdhttp.ErrServerClosed) {
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
