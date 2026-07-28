package admin

import (
	adminservice "github.com/anomalyco/hookah-store/notification-service/internal/services/admin"
	"github.com/anomalyco/hookah-store/notification-service/internal/transport/http"
	"github.com/gin-gonic/gin"
)

type AdminHandlers struct {
	adminService *adminservice.AdminService
}

func New(notifService *adminservice.AdminService) http.Handler {
	return &AdminHandlers{adminService: notifService}
}

func (h *AdminHandlers) Register(router *gin.RouterGroup) {
	notif := router.Group("/notifications")
	{
		notif.GET("user/:user_id", h.GetByUserID)
		notif.GET("email/:email", h.GetByEmail)
	}
}

func (h *AdminHandlers) ShutDown() {}

func (h *AdminHandlers) GetByUserID(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	notifications, err := h.adminService.GetByUserID(ctx, userID)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, notifications, "notifications retrieved successfully")
}

func (h *AdminHandlers) GetByEmail(ctx *gin.Context) {
	email := ctx.Param("email")

	notifications, err := h.adminService.GetByEmail(ctx, email)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, notifications, "notifications retrieved successfully")
}
