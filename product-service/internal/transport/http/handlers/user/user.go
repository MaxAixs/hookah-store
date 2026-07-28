package user

import (
	categoryservice "github.com/anomalyco/hookah-store/product-service/internal/services/category"
	productservice "github.com/anomalyco/hookah-store/product-service/internal/services/product"
	"github.com/anomalyco/hookah-store/product-service/internal/transport/http"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	categories *CategoryHandler
	products   *ProductHandler
}

func NewUserHandler(
	categoryService *categoryservice.CategoryService,
	productService *productservice.ProductService,
) http.PublicHandler {
	return &UserHandler{
		categories: &CategoryHandler{categoryService: categoryService},
		products:   &ProductHandler{productService: productService},
	}
}

func (h *UserHandler) RegisterPublic(router *gin.RouterGroup) {
	h.categories.registerRoutes(router)
	h.products.registerRoutes(router)
}
