package inventory

import (
	"context"
	"log/slog"
	"time"

	"github.com/anomalyco/hookah-store/product-service/internal/errs"
	"github.com/anomalyco/hookah-store/product-service/internal/models"
	"github.com/anomalyco/hookah-store/product-service/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type InventoryService struct {
	inventoryRepo repository.InventoryRepository
	validate      *validator.Validate
}

func New(inventoryRepo repository.InventoryRepository) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
		validate:      validator.New(),
	}
}

func (s *InventoryService) GetProductByID(ctx context.Context, productID uuid.UUID) (*models.InventoryResponse, error) {
	const fc = "product-service.services.GetProductByID"

	inventory, err := s.inventoryRepo.GetProductByID(ctx, productID)
	if err != nil {
		slog.Error("failed to get inventory", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	return models.InventoryToResponse(inventory), nil
}

func (s *InventoryService) GetAll(ctx context.Context) ([]*models.InventoryResponse, error) {
	const fc = "product-service.services.GetAllInventory"

	inventories, err := s.inventoryRepo.GetAllInventory(ctx)
	if err != nil {
		slog.Error("failed to get inventories", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	responses := make([]*models.InventoryResponse, 0, len(inventories))
	for i := range inventories {
		responses = append(responses, models.InventoryToResponse(&inventories[i]))
	}

	return responses, nil
}

func (s *InventoryService) UpdateInventoryProduct(ctx context.Context, productID uuid.UUID, req models.UpdateInventoryRequest) (*models.InventoryResponse, error) {
	const fc = "product-service.services.UpdateInventoryProduct"

	if err := s.validate.Struct(req); err != nil {
		slog.Error("validation failed", slog.String("fc", fc), slog.Any("error", err))

		return nil, err
	}

	inventory, err := s.inventoryRepo.GetProductByID(ctx, productID)
	if err != nil {
		slog.Error("failed to get inventory", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	inventory.Quantity = req.Quantity
	inventory.Reserved = req.Reserved
	inventory.UpdatedAt = time.Now()

	if err := s.inventoryRepo.UpdateProduct(ctx, inventory); err != nil {
		slog.Error("failed to update inventory", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	return models.InventoryToResponse(inventory), nil
}
