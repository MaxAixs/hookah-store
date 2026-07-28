package admin

import (
	"github.com/anomalyco/hookah-store/user-service/internal/errs"
	"github.com/anomalyco/hookah-store/user-service/internal/models"
	userservice "github.com/anomalyco/hookah-store/user-service/internal/services/user"
	"github.com/anomalyco/hookah-store/user-service/internal/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandlers struct {
	userService *userservice.UserService
}

func New(adminService *userservice.UserService) http.Handler {
	return &AdminHandlers{userService: adminService}
}

func (h *AdminHandlers) Register(router *gin.RouterGroup) {
	adminGroup := router.Group("/users")
	{
		adminGroup.POST("", h.CreateUser)
		adminGroup.PUT("/:id", h.UpdateUserByID)
		adminGroup.GET("/:id", h.GetUserByID)
		adminGroup.DELETE("/:id", h.DeleteUserByID)
	}
}

func (h *AdminHandlers) ShutDown() {}

const paramUserID = "id"

func (h *AdminHandlers) CreateUser(ctx *gin.Context) {
	var req models.CreateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.BadRequest(ctx, errs.ErrInvalidRequestBody)

		return
	}

	user, err := h.userService.CreateUser(ctx, req)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, user, "user created successfully")
}

func (h *AdminHandlers) UpdateUserByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param(paramUserID))
	if err != nil {
		http.BadRequest(ctx, errs.ErrInvalidUserID)

		return
	}

	var req models.UpdateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.BadRequest(ctx, errs.ErrInvalidRequestBody)

		return
	}

	user, err := h.userService.UpdateUserByID(ctx, id, req)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, user, "user updated successfully")
}

func (h *AdminHandlers) GetUserByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param(paramUserID))
	if err != nil {
		http.BadRequest(ctx, errs.ErrInvalidUserID)
	}

	user, err := h.userService.GetUserByID(ctx, id)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, user, "user retrieved successfully")
}

func (h *AdminHandlers) DeleteUserByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param(paramUserID))
	if err != nil {
		http.BadRequest(ctx, errs.ErrInvalidUserID)

		return
	}

	if err := h.userService.DeleteUser(ctx, id); err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, nil, "user deleted successfully")
}
