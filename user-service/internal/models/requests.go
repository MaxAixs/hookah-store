package models

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
}

type UpdateUserRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=8"`
	Role     string `json:"role" validate:"omitempty,oneof=admin user"`
}
