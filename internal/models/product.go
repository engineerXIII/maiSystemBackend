package models

import (
	"github.com/google/uuid"
	"time"
)

// Product full model
type Product struct {
	ProductID   uuid.UUID `json:"product_id" db:"product_id" validate:"omitempty"`
	Name        string    `json:"product_name" db:"product_name" validate:"required,lte=30"`
	Color       string    `json:"color" db:"color" validate:"required,lte=30"`
	Factory     string    `json:"factory" db:"factory" validate:"required,lte=30"`
	Description string    `json:"description" db:"description" validate:"required,lte=126"`
	Cost        int       `json:"cost" db:"cost" validate:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type ProductList struct {
	TotalCount int        `json:"total_count"`
	TotalPages int        `json:"total_pages"`
	Page       int        `json:"page"`
	Size       int        `json:"size"`
	HasMore    bool       `json:"has_more"`
	Products   []*Product `json:"products"`
}
