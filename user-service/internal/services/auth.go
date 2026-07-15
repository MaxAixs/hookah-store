package services

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/anomalyco/hookah-store/user-service/internal/errs"
	"github.com/anomalyco/hookah-store/user-service/internal/models"
	"github.com/anomalyco/hookah-store/user-service/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	jwtpkg "github.com/anomalyco/hookah-store/user-service/pkg/jwt"
)

type AuthService struct {
	db         *sqlx.DB
	outBoxRepo repository.OutBoxRepository
	userRepo   repository.UserRepository
	validate   *validator.Validate
	jwt        *jwtpkg.JwtConfig
}

func NewAuth(db *sqlx.DB, userRepo repository.UserRepository, jwtCfg *jwtpkg.JwtConfig) *AuthService {
	return &AuthService{
		db:       db,
		userRepo: userRepo,
		validate: validator.New(),
		jwt:      jwtCfg,
	}
}

const userRole = "user"
const userEventsTopic = "user.events"

func (s *AuthService) SignUp(ctx context.Context, req models.AuthRequest) error {
	const fc = "auth-service.services.CreateUser"

	if err := s.validate.Struct(req); err != nil {
		slog.Error("validation failed", slog.String("fc", fc), slog.Any("error", err))

		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		Role:         userRole,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("failed to create transaction", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}
	defer tx.Rollback()

	if err := s.userRepo.Create(ctx, tx, user); err != nil {
		slog.Error("failed to create user", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}

	if err := s.publishAuthEvent(ctx, tx, user.ID, user.Email, models.AuthEventSignUp, time.Now()); err != nil {
		slog.Error("failed to publish auth event", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}

	if err := tx.Commit(); err != nil {
		slog.Error("failed to commit", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}

	return nil
}

func (s *AuthService) SignIn(ctx context.Context, req models.AuthRequest) (string, error) {
	const fc = "auth-service.services.SignIn"

	if err := s.validate.Struct(req); err != nil {
		slog.Error("validation failed", slog.String("fc", fc), slog.Any("error", err))

		return "", err
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		slog.Error("failed to get user by email", slog.String("fc", fc), slog.Any("error", err))

		return "", errs.MapErr(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		slog.Error("Compare hash and password failed", slog.String("fc", fc), slog.Any("error", err))

		return "", errs.MapErr(err)
	}

	token, err := s.jwt.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		slog.Error("failed to generate token", slog.String("fc", fc), slog.Any("error", err))

		return "", errs.MapErr(err)
	}

	return token, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, req models.ResetPasswordRequest) error {
	const fc = "auth-service.services.ResetPassword"

	if err := s.validate.Struct(req.NewPassword); err != nil {
		slog.Error("validation failed", slog.String("fc", fc), slog.Any("error", err))

		return err
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		slog.Error("failed to get user", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("failed to create transaction", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}
	defer tx.Rollback()

	if err := s.userRepo.UpdatePassword(ctx, tx, req.Email, string(passwordHash)); err != nil {
		slog.Error("failed to update password", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}

	if err := s.publishAuthEvent(ctx, tx, user.ID, user.Email, models.AuthEventResetPassword, user.CreatedAt); err != nil {
		slog.Error("failed to publish auth event", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}

	if err := tx.Commit(); err != nil {
		slog.Error("failed to commit", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}

	return nil
}

func (s *AuthService) publishAuthEvent(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID, email string,
	eventType models.AuthEventType, timeStamp time.Time) error {
	const fc = "auth-service.services.publishAuthEvent"

	event := models.UserAuthEvent{
		UserID:    userID,
		Email:     email,
		EventType: eventType,
		TimeStamp: timeStamp,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal event", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}

	outBoxEvent := models.NewOutBoxEvent(userEventsTopic, userID.String(), eventType, payload)

	return s.outBoxRepo.SaveEvent(ctx, tx, outBoxEvent)
}
