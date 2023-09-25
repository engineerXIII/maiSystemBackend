package repository

import (
	"context"
	"encoding/json"
	"github.com/engineerXIII/maiSystemBackend/internal/models"
	"github.com/engineerXIII/maiSystemBackend/internal/order"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"time"
)

// Redis variables
const (
	basePrefix    = "api-orders:"
	cacheDuration = 3600
)

// Order redis repository
type orderRedisRepo struct {
	redisClient *redis.Client
}

// Order redis repository constructor
func NewOrderRedisRepo(redisClient *redis.Client) order.RedisRepository {
	return &orderRedisRepo{redisClient: redisClient}
}

func (n *orderRedisRepo) GetOrderKeysCtx(ctx context.Context) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderRedisRepo.GetOrderKeysCtx")
	defer span.Finish()
	var cursor uint64

	keys, cursor, err := n.redisClient.Scan(ctx, cursor, basePrefix+"*", 0).Result()
	if err != nil {
		return nil, errors.Wrap(err, "orderRedisRepo.GetOrderKeysCtx.redisClient.Scan")
	}
	return keys, nil
}

// Get order by id
func (n *orderRedisRepo) GetOrderByIDCtx(ctx context.Context, key string) (*models.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderRedisRepo.GetOrderByIDCtx")
	defer span.Finish()

	newsBytes, err := n.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "orderRedisRepo.GetOrderByIDCtx.redisClient.Get")
	}
	newsBase := &models.Order{}
	if err = json.Unmarshal(newsBytes, newsBase); err != nil {
		return nil, errors.Wrap(err, "orderRedisRepo.GetOrderByIDCtx.json.Unmarshal")
	}

	return newsBase, nil
}

// Cache order item
func (n *orderRedisRepo) SetOrderCtx(ctx context.Context, key string, seconds int, order *models.Order) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderRedisRepo.SetOrderCtx")
	defer span.Finish()

	orderBytes, err := json.Marshal(order)
	if err != nil {
		return errors.Wrap(err, "orderRedisRepo.SetOrderCtx.json.Marshal")
	}
	if err = n.redisClient.Set(ctx, key, orderBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return errors.Wrap(err, "orderRedisRepo.SetOrderCtx.redisClient.Set")
	}
	return nil
}

// Delete new item from cache
func (n *orderRedisRepo) DeleteOrderCtx(ctx context.Context, key string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderRedisRepo.DeleteOrderCtx")
	defer span.Finish()

	if err := n.redisClient.Del(ctx, key).Err(); err != nil {
		return errors.Wrap(err, "orderRedisRepo.DeleteOrderCtx.redisClient.Del")
	}
	return nil
}
