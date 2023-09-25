package usecase

import (
	"context"
	"encoding/json"
	"github.com/engineerXIII/maiSystemBackend/config"
	"github.com/engineerXIII/maiSystemBackend/internal/models"
	"github.com/engineerXIII/maiSystemBackend/internal/order"
	"github.com/engineerXIII/maiSystemBackend/pkg/httpErrors"
	"github.com/engineerXIII/maiSystemBackend/pkg/logger"
	"github.com/engineerXIII/maiSystemBackend/pkg/utils"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// Redis variables
const (
	basePrefix    = "api-orders:"
	cacheDuration = 3600
)

type orderUC struct {
	cfg       *config.Config
	orderRepo order.RedisRepository
	logger    logger.Logger
}

func NewOrderUseCase(cfg *config.Config, orderRepo order.RedisRepository, logger logger.Logger) order.UseCase {
	return &orderUC{cfg: cfg, orderRepo: orderRepo, logger: logger}
}

func (u *orderUC) Create(ctx context.Context, order *models.Order) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderUC.Create")
	defer span.Finish()

	if err := utils.ValidateStruct(ctx, order); err != nil {
		return httpErrors.NewBadRequestError(errors.WithMessage(err, "orderUC.Create.ValidateStruct"))
	}

	order.OrderId = uuid.New()
	order.Status = 1
	order.StatusMessage = order.Status.ToString()
	order.CalculateSum()

	s, _ := json.Marshal(order)
	u.logger.Debug(string(s))

	redisID := basePrefix + order.OrderId.String()

	err := u.orderRepo.SetOrderCtx(ctx, redisID, cacheDuration, order)
	if err != nil {
		return err
	}

	return err
}

func (u *orderUC) Update(ctx context.Context, order *models.Order) (*models.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderUC.Update")
	defer span.Finish()

	order.StatusMessage = order.Status.ToString()
	order.CalculateSum()

	redisID := basePrefix + order.OrderId.String()

	_, err := u.orderRepo.GetOrderByIDCtx(ctx, redisID)
	if err != nil {
		return nil, err
	}

	err = u.orderRepo.SetOrderCtx(ctx, redisID, cacheDuration, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (u *orderUC) GetOrderByID(ctx context.Context, orderUUID uuid.UUID) (*models.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderUC.GetOrderByID")
	defer span.Finish()

	redisID := basePrefix + orderUUID.String()

	p, err := u.orderRepo.GetOrderByIDCtx(ctx, redisID)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (u *orderUC) Delete(ctx context.Context, orderUUID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderUC.Delete")
	defer span.Finish()

	redisID := basePrefix + orderUUID.String()

	_, err := u.orderRepo.GetOrderByIDCtx(ctx, redisID)
	if err != nil {
		return err
	}

	if err := u.orderRepo.DeleteOrderCtx(ctx, redisID); err != nil {
		return err
	}

	return nil
}
