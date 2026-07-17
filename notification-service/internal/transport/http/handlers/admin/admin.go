package admin

import (
	adminservice "github.com/anomalyco/hookah-store/notification-service/internal/services/admin"
	"github.com/anomalyco/hookah-store/notification-service/internal/transport/http"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	adminService *adminservice.Service
}

func New(notifService *adminservice.Service) http.Handler {
	return &Handlers{adminService: notifService}
}

func (h *Handlers) Register(router *gin.RouterGroup) {
	adminGroup := router.Group("/admin/notifications")
	{
		adminGroup.GET("/user/:user_id", h.GetByUserID)
		adminGroup.GET("/email/:email", h.GetByEmail)
	}
}

func (h *Handlers) ShutDown() {}

func (h *Handlers) GetByUserID(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	notifications, err := h.adminService.GetByUserID(ctx, userID)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, notifications, "notifications retrieved successfully")
}

func (h *Handlers) GetByEmail(ctx *gin.Context) {
	email := ctx.Param("email")

	notifications, err := h.adminService.GetByEmail(ctx, email)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, notifications, "notifications retrieved successfully")
}
