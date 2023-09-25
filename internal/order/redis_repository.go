package order

import (
	"context"
	"github.com/engineerXIII/maiSystemBackend/internal/models"
)

type RedisRepository interface {
	GetOrderByIDCtx(ctx context.Context, key string) (*models.Order, error)
	SetOrderCtx(ctx context.Context, key string, seconds int, news *models.Order) error
	DeleteOrderCtx(ctx context.Context, key string) error
	GetOrderKeysCtx(ctx context.Context) ([]string, error)
}
