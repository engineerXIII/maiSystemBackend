//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package inventory

import (
	"context"
	"github.com/engineerXIII/maiSystemBackend/internal/models"
	"github.com/google/uuid"
)

type UseCase interface {
	AddItem(ctx context.Context, item *models.InventoryItem) (*models.InventoryItem, error)
	GetItemByID(ctx context.Context, item uuid.UUID) (*models.InventoryItem, error)
	RemoveItem(ctx context.Context, item *models.InventoryItem) (*models.InventoryItem, error)
}
