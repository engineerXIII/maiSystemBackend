package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/engineerXIII/maiSystemBackend/internal/inventory"
	"github.com/engineerXIII/maiSystemBackend/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"time"
)

const (
	basePrefix    = "api-inventory"
	cacheDuration = 3600
)

func NewInventoryRedisRepo(redisClient *redis.Client) inventory.RedisRepository {
	return &inventoryRedisRepo{redisClient: redisClient}
}

type inventoryRedisRepo struct {
	redisClient *redis.Client
}

func (i *inventoryRedisRepo) DeleteItemCtx(ctx context.Context, key string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "inventoryRedisRepo.DeleteItemCtx")
	defer span.Finish()

	if err := i.redisClient.Del(ctx, i.createKey(key)).Err(); err != nil {
		return errors.Wrap(err, "inventoryRedisRepo.DeleteItemCtx.redisClient.Del")
	}
	return nil
}

func (i *inventoryRedisRepo) GetByIDCtx(ctx context.Context, key string) (*models.InventoryItem, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "inventoryRedisRepo.GetByIDCtx")
	defer span.Finish()

	itemBytes, err := i.redisClient.Get(ctx, i.createKey(key)).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "authRedisRepo.GetByIDCtx.redisClient.Get")
	}
	item := &models.InventoryItem{}
	if err = json.Unmarshal(itemBytes, item); err != nil {
		return nil, errors.Wrap(err, "inventoryRedisRepo.GetByIDCtx.json.Unmarshal")
	}
	return item, nil
}

func (i *inventoryRedisRepo) SetItemCtx(ctx context.Context, key string, seconds int, item *models.InventoryItem) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "inventoryRedisRepo.SetItemCtx")
	defer span.Finish()

	itemBytes, err := json.Marshal(item)
	if err != nil {
		return errors.Wrap(err, "authRedisRepo.SetItemCtx.json.Unmarshal")
	}
	if err = i.redisClient.Set(ctx, i.createKey(key), itemBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return errors.Wrap(err, "inventoryRedisRepo.SetItemCtx.redisClient.Set")
	}
	return nil
}

func (i *inventoryRedisRepo) createKey(itemID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, itemID)
}
