//go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package order

import (
	"context"
	"github.com/engineerXIII/maiSystemBackend/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, product *models.Product) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) (*models.Product, error)
	GetProductByID(ctx context.Context, productID uuid.UUID) (*models.Product, error)
	Delete(ctx context.Context, productID uuid.UUID) error
}
