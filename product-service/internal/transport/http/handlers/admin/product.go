package admin

import (
	"github.com/anomalyco/hookah-store/product-service/internal/errs"
	"github.com/anomalyco/hookah-store/product-service/internal/models"
	productservice "github.com/anomalyco/hookah-store/product-service/internal/services/product"
	"github.com/anomalyco/hookah-store/product-service/internal/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandlers struct {
	productService *productservice.Service
}

func NewProductHandlers(productService *productservice.Service) http.Handler {
	return &ProductHandlers{productService: productService}
}

func (h *ProductHandlers) Register(router *gin.RouterGroup) {
	productGroup := router.Group("/products")
	{
		productGroup.POST("", h.CreateProduct)
		productGroup.GET("", h.GetAllProducts)
		productGroup.GET("/:id", h.GetProductByID)
		productGroup.PUT("/:id", h.UpdateProduct)
		productGroup.DELETE("/:id", h.DeleteProduct)
		productGroup.PUT("/:id/inventory", h.UpdateInventory)
	}
}

func (h *ProductHandlers) ShutDown() {}

const paramProductID = "id"

func (h *ProductHandlers) CreateProduct(ctx *gin.Context) {
	var req models.CreateProductRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.BadRequest(ctx, errs.ErrInvalidRequestBody)

		return
	}

	product, err := h.productService.CreateProduct(ctx, req)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, product, "product created successfully")
}

func (h *ProductHandlers) GetAllProducts(ctx *gin.Context) {
	products, err := h.productService.GetAllProducts(ctx)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, products, "products retrieved successfully")
}

func (h *ProductHandlers) GetProductByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param(paramProductID))
	if err != nil {
		http.BadRequest(ctx, errs.ErrInvalidProductID)

		return
	}

	product, err := h.productService.GetProductByID(ctx, id)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, product, "product retrieved successfully")
}

func (h *ProductHandlers) UpdateProduct(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param(paramProductID))
	if err != nil {
		http.BadRequest(ctx, errs.ErrInvalidProductID)

		return
	}

	var req models.UpdateProductRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.BadRequest(ctx, errs.ErrInvalidRequestBody)

		return
	}

	product, err := h.productService.UpdateProduct(ctx, id, req)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, product, "product updated successfully")
}

func (h *ProductHandlers) DeleteProduct(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param(paramProductID))
	if err != nil {
		http.BadRequest(ctx, errs.ErrInvalidProductID)

		return
	}

	if err := h.productService.DeleteProduct(ctx, id); err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, nil, "product deleted successfully")
}

func (h *ProductHandlers) UpdateInventory(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param(paramProductID))
	if err != nil {
		http.BadRequest(ctx, errs.ErrInvalidProductID)

		return
	}

	var req models.UpdateInventoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.BadRequest(ctx, errs.ErrInvalidRequestBody)

		return
	}

	inventory, err := h.productService.UpdateInventory(ctx, id, req)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, inventory, "inventory updated successfully")
}
