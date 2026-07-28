package product

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

func IsForeignKeyViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23503"
}

func (r *Repo) Create(ctx context.Context, tx *sqlx.Tx, product *models.Product) error {
	query := `
		INSERT INTO products (id, category_id, name, description, price, is_active, created_at, updated_at)
		VALUES (:id, :category_id, :name, :description, :price, :is_active, :created_at, :updated_at)`

	_, err := tx.NamedExecContext(ctx, query, product)
	if err != nil {
		if IsForeignKeyViolation(err) {
			return errs.ErrCategoryNotFound
		}

		return err
	}

	return nil
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	query := `
		SELECT id, category_id, name, COALESCE(description, '') AS description, price, is_active, created_at, updated_at
		FROM products
		WHERE id = $1`

	var product models.Product
	err := r.db.GetContext(ctx, &product, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrProductNotFound
		}

		return nil, err
	}

	return &product, nil
}

func (r *Repo) GetAll(ctx context.Context) ([]models.Product, error) {
	query := `
		SELECT id, category_id, name, COALESCE(description, '') AS description, price, is_active, created_at, updated_at
		FROM products
		ORDER BY created_at`

	var products []models.Product
	if err := r.db.SelectContext(ctx, &products, query); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Repo) Update(ctx context.Context, product *models.Product) error {
	query := `
		UPDATE products
		SET category_id = :category_id, name = :name, description = :description,
			price = :price, is_active = :is_active, updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, product)
	if err != nil {
		if IsForeignKeyViolation(err) {
			return errs.ErrCategoryNotFound
		}

		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errs.ErrProductNotFound
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errs.ErrProductNotFound
	}

	return nil
}

func (r *Repo) AddProduct(ctx context.Context, tx *sqlx.Tx, inventoryProduct *models.Inventory) error {
	query := `
		INSERT INTO inventory (product_id, quantity, reserved, updated_at)
		VALUES (:product_id, :quantity, :reserved, :updated_at)`

	_, err := tx.NamedExecContext(ctx, query, inventoryProduct)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetProductByID(ctx context.Context, productID uuid.UUID) (*models.Inventory, error) {
	query := `
		SELECT product_id, quantity, reserved, updated_at
		FROM inventory
		WHERE product_id = $1`

	var inventory models.Inventory
	err := r.db.GetContext(ctx, &inventory, query, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrProductNotFound
		}

		return nil, err
	}

	return &inventory, nil
}

func (r *Repo) GetAllInventory(ctx context.Context) ([]models.Inventory, error) {
	query := `
		SELECT product_id, quantity, reserved, updated_at
		FROM inventory
		ORDER BY updated_at`

	var inventories []models.Inventory
	if err := r.db.SelectContext(ctx, &inventories, query); err != nil {
		return nil, err
	}

	return inventories, nil
}

func (r *Repo) UpdateProduct(ctx context.Context, inventory *models.Inventory) error {
	query := `
		UPDATE inventory
		SET quantity = :quantity, reserved = :reserved, updated_at = :updated_at
		WHERE product_id = :product_id`

	result, err := r.db.NamedExecContext(ctx, query, inventory)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errs.ErrProductNotFound
	}

	return nil
}
