package admin

import (
	"github.com/anomalyco/hookah-store/product-service/internal/errs"
	"github.com/anomalyco/hookah-store/product-service/internal/models"
	categoryservice "github.com/anomalyco/hookah-store/product-service/internal/services/category"
	"github.com/anomalyco/hookah-store/product-service/internal/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandlers struct {
	categoryService *categoryservice.Service
}

func NewCategoryHandlers(categoryService *categoryservice.Service) http.Handler {
	return &CategoryHandlers{categoryService: categoryService}
}

func (h *CategoryHandlers) Register(router *gin.RouterGroup) {
	categoryGroup := router.Group("/categories")
	{
		categoryGroup.POST("", h.CreateCategory)
		categoryGroup.GET("", h.GetAllCategories)
		categoryGroup.GET("/:id", h.GetCategoryByID)
		categoryGroup.PUT("/:id", h.UpdateCategory)
		categoryGroup.DELETE("/:id", h.DeleteCategory)
	}
}

func (h *CategoryHandlers) ShutDown() {}

const paramCategoryID = "id"

func (h *CategoryHandlers) CreateCategory(ctx *gin.Context) {
	var req models.CreateCategoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.BadRequest(ctx, errs.ErrInvalidRequestBody)

		return
	}

	category, err := h.categoryService.CreateCategory(ctx, req)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, category, "category created successfully")
}

func (h *CategoryHandlers) GetAllCategories(ctx *gin.Context) {
	categories, err := h.categoryService.GetAllCategories(ctx)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, categories, "categories retrieved successfully")
}

func (h *CategoryHandlers) GetCategoryByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param(paramCategoryID))
	if err != nil {
		http.BadRequest(ctx, errs.ErrInvalidCategoryID)

		return
	}

	category, err := h.categoryService.GetCategoryByID(ctx, id)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, category, "category retrieved successfully")
}

func (h *CategoryHandlers) UpdateCategory(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param(paramCategoryID))
	if err != nil {
		http.BadRequest(ctx, errs.ErrInvalidCategoryID)

		return
	}

	var req models.UpdateCategoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.BadRequest(ctx, errs.ErrInvalidRequestBody)

		return
	}

	category, err := h.categoryService.UpdateCategory(ctx, id, req)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, category, "category updated successfully")
}

func (h *CategoryHandlers) DeleteCategory(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param(paramCategoryID))
	if err != nil {
		http.BadRequest(ctx, errs.ErrInvalidCategoryID)

		return
	}

	if err := h.categoryService.DeleteCategory(ctx, id); err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, nil, "category deleted successfully")
}
