package auth

import (
	"github.com/anomalyco/hookah-store/user-service/internal/errs"
	"github.com/anomalyco/hookah-store/user-service/internal/models"
	authservice "github.com/anomalyco/hookah-store/user-service/internal/services/auth"
	"github.com/anomalyco/hookah-store/user-service/internal/transport/http"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	authService *authservice.Service
}

func New(authService *authservice.Service) http.Handler {
	return &Handlers{
		authService: authService,
	}
}

func (h *Handlers) Register(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/sign-up", h.SignUp)
		authGroup.POST("/sign-in", h.SignIn)
		authGroup.POST("/reset-password", h.ResetPassword)
	}
}

func (h *Handlers) ShutDown() {}

func (h *Handlers) SignUp(ctx *gin.Context) {
	var req models.AuthRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.BadRequest(ctx, errs.ErrInvalidRequestBody)

		return
	}

	if err := h.authService.SignUp(ctx, req); err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, "user successfully signed up")
}

func (h *Handlers) SignIn(ctx *gin.Context) {
	var req models.AuthRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.BadRequest(ctx, errs.ErrInvalidRequestBody)

		return
	}

	token, err := h.authService.SignIn(ctx, req)
	if err != nil {
		http.BadRequest(ctx, err)

		return
	}

	http.OK(ctx, token, "user successfully signed in")
}

func (h *Handlers) ResetPassword(ctx *gin.Context) {
	var req models.ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.BadRequest(ctx, errs.ErrInvalidRequestBody)

		return
	}

	if err := h.authService.ResetPassword(ctx, req); err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, "user successfully reset password")
}
