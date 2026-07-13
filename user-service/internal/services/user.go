package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/anomalyco/hookah-store/user-service/internal/errs"
	"github.com/anomalyco/hookah-store/user-service/internal/models"
	"github.com/anomalyco/hookah-store/user-service/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.UserRepository
	validate *validator.Validate
}

func NewAdmin(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		validate: validator.New(),
	}
}

func (s *UserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error) {
	const fc = "user-service.services.CreateUser"

	if err := s.validate.Struct(req); err != nil {
		slog.Error("validation failed", slog.String("fc", fc), slog.Any("error", err))

		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		Role:         req.Role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, nil, user); err != nil {
		slog.Error("failed to create user", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	userResponse := models.UserToResponse(user)

	return userResponse, nil
}

func (s *UserService) UpdateUserByID(ctx context.Context, id uuid.UUID, req models.UpdateUserRequest) (*models.UserResponse, error) {
	const fc = "user-service.services.UpdateUser"

	if err := s.validate.Struct(req); err != nil {
		slog.Error("validation failed", slog.String("fc", fc), slog.Any("error", err))

		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("failed to get user", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("password hash failed", slog.String("fc", fc), slog.Any("error", err))

			return nil, errs.MapErr(err)
		}
		user.PasswordHash = string(hash)
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		slog.Error("failed to update user", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	userResponse := models.UserToResponse(user)

	return userResponse, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserResponse, error) {
	const fc = "user-service.services.GetUserByID"

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("failed to get user", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	userResponse := models.UserToResponse(user)

	return userResponse, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	const fc = "user-service.services.DeleteUser"

	if err := s.userRepo.Delete(ctx, id); err != nil {
		slog.Error("failed to delete user", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}

	return nil
}
