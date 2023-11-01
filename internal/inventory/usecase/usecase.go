package usecase

import (
	"context"
	"github.com/engineerXIII/maiSystemBackend/config"
	"github.com/engineerXIII/maiSystemBackend/internal/inventory"
	"github.com/engineerXIII/maiSystemBackend/internal/models"
	"github.com/engineerXIII/maiSystemBackend/pkg/logger"
	"github.com/google/uuid"
)

type inventoryUC struct {
	cfg       *config.Config
	redisRepo inventory.RedisRepository
	logger    logger.Logger
}

func NewInventoryUseCase(cfg *config.Config, redisRepo inventory.RedisRepository, log logger.Logger) inventory.UseCase {
	return &inventoryUC{cfg: cfg, redisRepo: redisRepo, logger: log}
}

func (i *inventoryUC) AddItem(ctx context.Context, item *models.InventoryItem) (*models.InventoryItem, error) {
	inventoryItem, err := i.redisRepo.GetByIDCtx(ctx, item.UUID.String())
	if err != nil {
		return nil, err
	}
	inventoryItem.Qty = inventoryItem.Qty + 1
	return nil, i.redisRepo.SetItemCtx(ctx, inventoryItem.UUID.String(), 3600, inventoryItem)
}

func (i *inventoryUC) GetItemByID(ctx context.Context, item uuid.UUID) (*models.InventoryItem, error) {
	return i.GetItemByID(ctx, item)
}

func (i *inventoryUC) RemoveItem(ctx context.Context, item *models.InventoryItem) (*models.InventoryItem, error) {
	inventoryItem, err := i.redisRepo.GetByIDCtx(ctx, item.UUID.String())
	if err != nil {
		return nil, err
	}
	inventoryItem.Qty = inventoryItem.Qty - 1
	if inventoryItem.Qty == 0 {
		return nil, i.redisRepo.DeleteItemCtx(ctx, item.UUID.String())
	}
	return nil, i.redisRepo.SetItemCtx(ctx, inventoryItem.UUID.String(), 3600, inventoryItem)
}
