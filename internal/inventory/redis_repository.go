//go:generate mockgen -source redis_repository.go -destination mock/redis_repository_mock.go -package mock
package inventory

import (
	"context"
	"github.com/engineerXIII/maiSystemBackend/internal/models"
)

type RedisRepository interface {
	GetByIDCtx(ctx context.Context, key string) (*models.InventoryItem, error)
	SetItemCtx(ctx context.Context, key string, seconds int, item *models.InventoryItem) error
	DeleteItemCtx(ctx context.Context, key string) error
}
