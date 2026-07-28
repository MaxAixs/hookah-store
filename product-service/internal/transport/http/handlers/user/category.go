package user

import (
	"github.com/anomalyco/hookah-store/product-service/internal/errs"
	categoryservice "github.com/anomalyco/hookah-store/product-service/internal/services/category"
	"github.com/anomalyco/hookah-store/product-service/internal/transport/http"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService *categoryservice.CategoryService
}

func (h *CategoryHandler) registerRoutes(router *gin.RouterGroup) {
	categories := router.Group("/categories")
	{
		categories.GET("/:slug", h.GetCategoryBySlug)
	}
}

const paramCategorySlug = "slug"

func (h *CategoryHandler) GetCategoryBySlug(ctx *gin.Context) {
	slug := ctx.Param(paramCategorySlug)
	if slug == "" {
		http.BadRequest(ctx, errs.ErrInvalidRequestBody)

		return
	}

	category, err := h.categoryService.GetCategoryBySlug(ctx, slug)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, category, "category retrieved successfully")
}
