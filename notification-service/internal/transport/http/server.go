package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/anomalyco/hookah-store/notification-service/internal/config"
	"github.com/anomalyco/hookah-store/user-service/pkg/auth"
	jwtpkg "github.com/anomalyco/hookah-store/user-service/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type Server struct {
	srv    *http.Server
	router *gin.Engine
}

func New(cfg *config.HTTPServerConfig, jwtCfg *jwtpkg.JwtConfig, webhookHandlers Handler, adminHandlers Handler) *Server {
	router := gin.New()

	api := router.Group("/api")
	adminApi := router.Group("/api/admin", auth.RequireAdminRole(jwtCfg))
	webhookHandlers.Register(api)
	adminHandlers.Register(adminApi)

	return &Server{
		srv: &http.Server{
			Addr:              fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Handler:           router,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
		},
		router: router,
	}
}

func (s *Server) Run() error {
	slog.Info("server listening on", slog.String("addr", s.srv.Addr))
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("server shutting down")
	return s.srv.Shutdown(ctx)
}
