package models

import "github.com/google/uuid"

type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Slug        string `json:"slug" validate:"required,min=2,max=255"`
	Description string `json:"description" validate:"omitempty"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2,max=255"`
	Slug        string `json:"slug" validate:"omitempty,min=2,max=255"`
	Description string `json:"description" validate:"omitempty"`
}

type CreateProductRequest struct {
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=2,max=255"`
	Description string    `json:"description" validate:"omitempty"`
	Price       float64   `json:"price" validate:"required,gt=0"`
	Stock       int       `json:"stock" validate:"gte=0"`
}

type UpdateProductRequest struct {
	CategoryID  uuid.UUID `json:"category_id" validate:"omitempty"`
	Name        string    `json:"name" validate:"omitempty,min=2,max=255"`
	Description string    `json:"description" validate:"omitempty"`
	Price       float64   `json:"price" validate:"omitempty,gt=0"`
	IsActive    *bool     `json:"is_active" validate:"omitempty"`
}

type UpdateInventoryRequest struct {
	Quantity int `json:"quantity" validate:"gte=0"`
	Reserved int `json:"reserved" validate:"gte=0"`
}
