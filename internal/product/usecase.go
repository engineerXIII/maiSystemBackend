//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package product

import (
	"context"
	"github.com/engineerXIII/maiSystemBackend/internal/models"
	"github.com/engineerXIII/maiSystemBackend/pkg/utils"
	"github.com/google/uuid"
)

// Product use case
type UseCase interface {
	Create(ctx context.Context, product *models.Product) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) (*models.Product, error)
	GetProductByID(ctx context.Context, productID uuid.UUID) (*models.Product, error)
	Delete(ctx context.Context, productID uuid.UUID) error
	GetProducts(ctx context.Context, pq *utils.PaginationQuery) (*models.ProductList, error)
	SearchByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.ProductList, error)
}
