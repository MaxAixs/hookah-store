package repository

import (
	"context"

	"github.com/anomalyco/hookah-store/product-service/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	GetBySlug(ctx context.Context, slug string) (*models.Category, error)
	GetAll(ctx context.Context) ([]models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProductRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, product *models.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	GetAll(ctx context.Context) ([]models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type InventoryRepository interface {
	AddProduct(ctx context.Context, tx *sqlx.Tx, inventoryProduct *models.Inventory) error
	GetProductByID(ctx context.Context, productID uuid.UUID) (*models.Inventory, error)
	GetAllInventory(ctx context.Context) ([]models.Inventory, error)
	UpdateProduct(ctx context.Context, inventoryProduct *models.Inventory) error
}
