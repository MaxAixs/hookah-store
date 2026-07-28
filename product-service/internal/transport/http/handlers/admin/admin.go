package admin

import (
	"github.com/anomalyco/hookah-store/product-service/internal/services/category"
	"github.com/anomalyco/hookah-store/product-service/internal/services/inventory"
	"github.com/anomalyco/hookah-store/product-service/internal/services/product"
	"github.com/anomalyco/hookah-store/product-service/internal/transport/http"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	categories *CategoryHandlers
	products   *ProductHandlers
	inventory  *InventoryHandlers
}

func New(
	categoryService *category.CategoryService,
	productService *product.ProductService,
	inventoryService *inventory.InventoryService,
) http.Handler {
	return &AdminHandler{
		categories: &CategoryHandlers{categoryService: categoryService},
		products:   &ProductHandlers{productService: productService},
		inventory:  &InventoryHandlers{inventoryService: inventoryService},
	}
}

func (h *AdminHandler) Register(router *gin.RouterGroup) {
	h.categories.registerRoutes(router)
	h.products.registerRoutes(router)
	h.inventory.registerRoutes(router)
}

func (h *AdminHandler) ShutDown() {}
