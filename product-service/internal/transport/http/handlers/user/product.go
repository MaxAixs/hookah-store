package user

import (
	productservice "github.com/anomalyco/hookah-store/product-service/internal/services/product"
	"github.com/anomalyco/hookah-store/product-service/internal/transport/http"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *productservice.ProductService
}

func NewProductHandler(productService *productservice.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) RegisterPublic(router *gin.RouterGroup) {
	products := router.Group("/products")
	{
		products.GET("", h.GetAllProducts)
	}
}

func (h *ProductHandler) GetAllProducts(ctx *gin.Context) {
	products, err := h.productService.GetAllProducts(ctx)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, products, "products retrieved successfully")
}
