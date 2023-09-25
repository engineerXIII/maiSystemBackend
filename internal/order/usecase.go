//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package order

import (
	"context"
	"github.com/engineerXIII/maiSystemBackend/internal/models"
	"github.com/google/uuid"
)

// Product use case
type UseCase interface {
	Create(ctx context.Context, order *models.Order) (*models.Order, error)
	Update(ctx context.Context, order *models.Order) (*models.Order, error)
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (*models.Order, error)
	Delete(ctx context.Context, orderID uuid.UUID) error
}
