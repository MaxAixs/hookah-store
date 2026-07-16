package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/anomalyco/hookah-store/user-service/internal/errs"
	"github.com/anomalyco/hookah-store/user-service/internal/models"
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

func (r *Repo) Create(ctx context.Context, tx *sqlx.Tx, user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, role, created_at, updated_at)
		VALUES (:id, :email, :password_hash, :role, :created_at, :updated_at)`

	_, err := tx.NamedExecContext(ctx, query, user)
	if err != nil {
		if IsEmailAlreadyExist(err) {
			return errs.ErrUserAlreadyExists
		}

		return err
	}

	return nil
}

func IsEmailAlreadyExist(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = $1`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repo) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET email = :email, password_hash = :password_hash, role = :role, updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}

func (r *Repo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repo) UpdatePassword(ctx context.Context, tx *sqlx.Tx, email string, newPassword string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE email = $2`

	result, err := tx.ExecContext(ctx, query, newPassword, email)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}
