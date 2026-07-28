package category

import (
	"context"
	"database/sql"
	"errors"

	"github.com/anomalyco/hookah-store/product-service/internal/errs"
	"github.com/anomalyco/hookah-store/product-service/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, category *models.Category) error {
	query := `
		INSERT INTO categories (id, name, slug, description, created_at, updated_at)
		VALUES (:id, :name, :slug, :description, :created_at, :updated_at)`

	_, err := r.db.NamedExecContext(ctx, query, category)
	if err != nil {
		if IsSlugAlreadyExist(err) {
			return errs.ErrCategoryAlreadyExist
		}

		return err
	}

	return nil
}

func IsSlugAlreadyExist(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}

func IsForeignKeyViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23503"
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	query := `
		SELECT id, name, slug, COALESCE(description, '') AS description, created_at, updated_at
		FROM categories
		WHERE id = $1`

	var category models.Category
	err := r.db.GetContext(ctx, &category, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCategoryNotFound
		}

		return nil, err
	}

	return &category, nil
}

func (r *Repo) GetAll(ctx context.Context) ([]models.Category, error) {
	query := `
		SELECT id, name, slug, COALESCE(description, '') AS description, created_at, updated_at
		FROM categories
		ORDER BY created_at`

	var categories []models.Category
	if err := r.db.SelectContext(ctx, &categories, query); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *Repo) Update(ctx context.Context, category *models.Category) error {
	query := `
		UPDATE categories
		SET name = :name, slug = :slug, description = :description, updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, category)
	if err != nil {
		if IsSlugAlreadyExist(err) {
			return errs.ErrCategoryAlreadyExist
		}

		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errs.ErrCategoryNotFound
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		if IsForeignKeyViolation(err) {
			return errs.ErrCategoryHasProducts
		}

		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errs.ErrCategoryNotFound
	}

	return nil
}
