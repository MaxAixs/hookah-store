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
	emailservice "github.com/anomalyco/hookah-store/notification-service/internal/services/email"
	"github.com/anomalyco/hookah-store/notification-service/internal/transport/http"
	adminHandlers "github.com/anomalyco/hookah-store/notification-service/internal/transport/http/handlers/admin"
	"github.com/anomalyco/hookah-store/notification-service/internal/transport/http/handlers/consumer"
	webhookHandlers "github.com/anomalyco/hookah-store/notification-service/internal/transport/http/handlers/webhook"
	"github.com/anomalyco/hookah-store/notification-service/pkg/database"
	kafkapkg "github.com/anomalyco/hookah-store/notification-service/pkg/kafka"
	mailgunpkg "github.com/anomalyco/hookah-store/notification-service/pkg/mailgun"
	jwtpkg "github.com/anomalyco/hookah-store/user-service/pkg/jwt"
)

const (
	serviceName     = "notification-service"
	userEventsTopic = "user.events"
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

	notifRepo := notification.New(db)

	mailgunClient := mailgunpkg.New(cfg.MailGun)

	adminService := admin.New(notifRepo)
	emailService := emailservice.New(notifRepo, cfg.MailGun.WebhookSigningKey, mailgunClient)

	jwtCfg := jwtpkg.New(cfg.JWT.Secret, cfg.JWT.TTL)

	emailHandler := consumer.New(emailService, userEventsTopic)
	kafkaConsumer := kafkapkg.New(cfg.Kafka)
	emailHandler.Register(kafkaConsumer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		slog.Info("starting kafka consumer", slog.String("group_id", cfg.Kafka.GroupID))
		if err := kafkaConsumer.Start(ctx); err != nil {
			slog.Error("kafka consumer failed", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	adminHandler := adminHandlers.New(adminService)
	webhookHandler := webhookHandlers.New(emailService)

	httpServer := http.New(&cfg.HTTPServer, jwtCfg, webhookHandler, adminHandler)
	go func() {
		slog.Info("starting http server", slog.String("port", cfg.HTTPServer.Port))
		if err := httpServer.Run(); err != nil {
			slog.Error("failed to start http server", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	slog.Info("received shutdown signal", slog.String("signal", sig.String()))

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown http server", slog.String("err", err.Error()))
	}

	if err := kafkaConsumer.Close(); err != nil {
		slog.Error("failed to close kafka consumer", slog.String("err", err.Error()))
	}

	slog.Info("service stopped", slog.String("service", serviceName))
}
