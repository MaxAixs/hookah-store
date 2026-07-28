package category

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

type CategoryService struct {
	categoryRepo repository.CategoryRepository
	validate     *validator.Validate
}

func New(categoryRepo repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		validate:     validator.New(),
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, req models.CreateCategoryRequest) (*models.CategoryResponse, error) {
	const fc = "product-service.services.CreateCategory"

	if err := s.validate.Struct(req); err != nil {
		slog.Error("validation failed", slog.String("fc", fc), slog.Any("error", err))

		return nil, err
	}

	category := &models.Category{
		ID:          uuid.New(),
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		slog.Error("failed to create category", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	return models.CategoryToResponse(category), nil
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.CategoryResponse, error) {
	const fc = "product-service.services.GetCategoryByID"

	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("failed to get category", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	return models.CategoryToResponse(category), nil
}

func (s *CategoryService) GetCategoryBySlug(ctx context.Context, slug string) (*models.CategoryResponse, error) {
	const fc = "product-service.services.GetCategoryBySlug"

	category, err := s.categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		slog.Error("failed to get category by slug", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	return models.CategoryToResponse(category), nil
}

func (s *CategoryService) GetAllCategories(ctx context.Context) ([]*models.CategoryResponse, error) {
	const fc = "product-service.services.GetAllCategories"

	categories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {
		slog.Error("failed to get categories", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	responses := make([]*models.CategoryResponse, 0, len(categories))
	for i := range categories {
		responses = append(responses, models.CategoryToResponse(&categories[i]))
	}

	return responses, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id uuid.UUID, req models.UpdateCategoryRequest) (*models.CategoryResponse, error) {
	const fc = "product-service.services.UpdateCategory"

	if err := s.validate.Struct(req); err != nil {
		slog.Error("validation failed", slog.String("fc", fc), slog.Any("error", err))

		return nil, err
	}

	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("failed to get category", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	applyCategoryUpdates(category, req)

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		slog.Error("failed to update category", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	return models.CategoryToResponse(category), nil
}

func applyCategoryUpdates(category *models.Category, req models.UpdateCategoryRequest) {
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		category.Slug = req.Slug
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	category.UpdatedAt = time.Now()
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	const fc = "product-service.services.DeleteCategory"

	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		slog.Error("failed to delete category", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}

	return nil
}
