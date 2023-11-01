package grpc

import (
	"context"
	"github.com/engineerXIII/maiSystemBackend/config"
	"github.com/engineerXIII/maiSystemBackend/internal/inventory"
	"github.com/engineerXIII/maiSystemBackend/pkg/logger"
	pb "github.com/engineerXIII/maiSystemBackend/proto/api/v1"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

type InventoryServer struct {
	cfg         *config.Config
	inventoryUC inventory.UseCase
	logger      logger.Logger
	pb.UnimplementedInventoryServiceServer
}

func NewInventoryServer(cfg *config.Config, inventoryUC inventory.UseCase, logger logger.Logger) *InventoryServer {
	return &InventoryServer{
		inventoryUC: inventoryUC,
		logger:      logger,
		cfg:         cfg,
	}
}

func (s InventoryServer) CheckItem(ctx context.Context, in *pb.ItemRequest) (*pb.ItemAvailableResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "inventory.CheckItem")
	defer span.Finish()

	response := &pb.ItemAvailableResponse{}
	for _, item := range in.Item {
		uuidItem, _ := uuid.FromBytes([]byte(item.Uuid))
		foundItem, err := s.inventoryUC.GetItemByID(ctx, uuidItem)
		if err != nil {
			return nil, err
		} else if foundItem == nil {
			return &pb.ItemAvailableResponse{
				Item:   []*pb.Item{item},
				Status: pb.Status_NotFound,
			}, nil
		} else {
			response.Item = append(response.Item, &pb.Item{
				Uuid: foundItem.UUID.String(),
				Qty:  uint64(item.Qty),
			})
		}
	}
	response.Status = pb.Status_OK
	return response, nil
}

func (s InventoryServer) AddItem(ctx context.Context, in *pb.ItemRequest) (*pb.Response, error) {
	return nil, nil
}

func (s InventoryServer) RemoveItem(ctx context.Context, in *pb.ItemRequest) (*pb.Response, error) {
	return nil, nil
}
