package product

import (
	"context"
	"log/slog"
	"time"

	"github.com/anomalyco/hookah-store/product-service/internal/errs"
	"github.com/anomalyco/hookah-store/product-service/internal/models"
	"github.com/anomalyco/hookah-store/product-service/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProductService struct {
	db            *sqlx.DB
	productRepo   repository.ProductRepository
	inventoryRepo repository.InventoryRepository
	validate      *validator.Validate
}

func New(db *sqlx.DB, productRepo repository.ProductRepository, inventoryRepo repository.InventoryRepository) *ProductService {
	return &ProductService{
		db:            db,
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
		validate:      validator.New(),
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, req models.CreateProductRequest) (*models.ProductResponse, error) {
	const fc = "product-service.services.CreateProduct"

	if err := s.validate.Struct(req); err != nil {
		slog.Error("validation failed", slog.String("fc", fc), slog.Any("error", err))

		return nil, err
	}

	product := &models.Product{
		ID:          uuid.New(),
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	inventoryProduct := &models.Inventory{
		ProductID: product.ID,
		Quantity:  req.Stock,
		Reserved:  0,
		UpdatedAt: time.Now(),
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("failed to create transaction", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}
	defer tx.Rollback()

	if err := s.productRepo.Create(ctx, tx, product); err != nil {
		slog.Error("failed to create product", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	if err := s.inventoryRepo.AddProduct(ctx, tx, inventoryProduct); err != nil {
		slog.Error("failed to create inventory", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	if err := tx.Commit(); err != nil {
		slog.Error("failed to commit", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	return models.ProductToResponse(product), nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id uuid.UUID) (*models.ProductResponse, error) {
	const fc = "product-service.services.GetProductByID"

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("failed to get product", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	return models.ProductToResponse(product), nil
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]*models.ProductResponse, error) {
	const fc = "product-service.services.GetAllProducts"

	products, err := s.productRepo.GetAll(ctx)
	if err != nil {
		slog.Error("failed to get products", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	responses := make([]*models.ProductResponse, 0, len(products))
	for i := range products {
		responses = append(responses, models.ProductToResponse(&products[i]))
	}

	return responses, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id uuid.UUID, req models.UpdateProductRequest) (*models.ProductResponse, error) {
	const fc = "product-service.services.UpdateProduct"

	if err := s.validate.Struct(req); err != nil {
		slog.Error("validation failed", slog.String("fc", fc), slog.Any("error", err))

		return nil, err
	}

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("failed to get product", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	applyProductUpdates(product, req)

	if err := s.productRepo.Update(ctx, product); err != nil {
		slog.Error("failed to update product", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	return models.ProductToResponse(product), nil
}

func applyProductUpdates(product *models.Product, req models.UpdateProductRequest) {
	if req.CategoryID != uuid.Nil {
		product.CategoryID = req.CategoryID
	}
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != 0 {
		product.Price = req.Price
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	product.UpdatedAt = time.Now()
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	const fc = "product-service.services.DeleteProduct"

	if err := s.productRepo.Delete(ctx, id); err != nil {
		slog.Error("failed to delete product", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}

	return nil
}
