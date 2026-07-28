package admin

import (
	"github.com/anomalyco/hookah-store/product-service/internal/errs"
	"github.com/anomalyco/hookah-store/product-service/internal/models"
	inventoryservice "github.com/anomalyco/hookah-store/product-service/internal/services/inventory"
	"github.com/anomalyco/hookah-store/product-service/internal/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InventoryHandlers struct {
	inventoryService *inventoryservice.InventoryService
}

func (h *InventoryHandlers) registerRoutes(router *gin.RouterGroup) {
	inventoryGroup := router.Group("/inventory")
	{
		inventoryGroup.GET("", h.GetAll)
		inventoryGroup.GET("/:id", h.GetProductByID)
		inventoryGroup.PUT("/:id", h.UpdateInventoryProduct)
	}
}

const paramInventoryProductID = "id"

func (h *InventoryHandlers) GetAll(ctx *gin.Context) {
	inventories, err := h.inventoryService.GetAll(ctx)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, inventories, "inventory retrieved successfully")
}

func (h *InventoryHandlers) GetProductByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param(paramInventoryProductID))
	if err != nil {
		http.BadRequest(ctx, errs.ErrInvalidProductID)

		return
	}

	inventory, err := h.inventoryService.GetProductByID(ctx, id)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, inventory, "inventory item retrieved successfully")
}

func (h *InventoryHandlers) UpdateInventoryProduct(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param(paramInventoryProductID))
	if err != nil {
		http.BadRequest(ctx, errs.ErrInvalidProductID)

		return
	}

	var req models.UpdateInventoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.BadRequest(ctx, errs.ErrInvalidRequestBody)

		return
	}

	inventory, err := h.inventoryService.UpdateInventoryProduct(ctx, id, req)
	if err != nil {
		http.HandleServiceError(ctx, err)

		return
	}

	http.OK(ctx, inventory, "inventory item updated successfully")
}
