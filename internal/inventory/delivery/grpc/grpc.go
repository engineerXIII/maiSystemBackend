package grpc

import (
	"context"
	"github.com/engineerXIII/maiSystemBackend/config"
	"github.com/engineerXIII/maiSystemBackend/internal/inventory"
	"github.com/engineerXIII/maiSystemBackend/internal/models"
	"github.com/engineerXIII/maiSystemBackend/pkg/logger"
	pb "github.com/engineerXIII/maiSystemBackend/proto/api/v1"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

type InventoryServer struct {
	cfg         *config.Config
	inventoryUC inventory.UseCase
	logger      logger.Logger
	pb.InventoryServiceServer
}

func NewInventoryServer(cfg *config.Config, inventoryUC inventory.UseCase, logger logger.Logger) *InventoryServer {
	return &InventoryServer{
		inventoryUC: inventoryUC,
		logger:      logger,
		cfg:         cfg,
	}
}

func (s InventoryServer) CheckItem(c context.Context, in *pb.ItemRequest) (*pb.ItemAvailableResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "inventory.CheckItem")
	defer span.Finish()

	response := &pb.ItemAvailableResponse{}
	status := pb.Status_OK
	for _, item := range in.Item {
		uuidItem, _ := uuid.Parse(item.Uuid)
		foundItem, err := s.inventoryUC.GetItemByID(ctx, uuidItem)
		if err != nil {
			return nil, err
		} else if foundItem == nil {
			response.Items = append(response.Items, &pb.ItemAvailableStatus{
				Item: &pb.Item{
					Uuid: item.Uuid,
					Qty:  0,
				},
				Status: pb.Status_NotFound,
			})
			status = pb.Status_NotEnoughAvailable
		} else {
			qty := item.Qty
			sts := pb.Status_OK
			if uint64(foundItem.Qty) < item.Qty {
				qty = uint64(foundItem.Qty)
				sts = pb.Status_NotEnoughAvailable
				status = pb.Status_NotEnoughAvailable
			}
			response.Items = append(response.Items, &pb.ItemAvailableStatus{
				Item: &pb.Item{
					Uuid: foundItem.UUID.String(),
					Qty:  qty,
				},
				Status: sts,
			})
		}
	}
	response.Status = status
	return response, nil
}

func (s InventoryServer) AddItem(c context.Context, in *pb.ItemRequest) (*pb.Response, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "inventory.CheckItem")
	defer span.Finish()

	response := &pb.Response{}
	uid, _ := uuid.Parse(in.Item[0].Uuid)
	item := &models.InventoryItem{
		UUID: uid,
		Qty:  int(in.Item[0].Qty),
	}
	_, err := s.inventoryUC.AddItem(ctx, item)
	if err != nil {
		return nil, err
	}
	response.StatusMessage = "Created"
	response.Status = pb.Status_OK
	return response, nil
}

func (s InventoryServer) RemoveItem(c context.Context, in *pb.ItemRequest) (*pb.Response, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "inventory.CheckItem")
	defer span.Finish()

	response := &pb.Response{}
	status := pb.Status_OK
	response.Status = status
	response.StatusMessage = "Successfull"
	for _, item := range in.Item {
		uuidItem, _ := uuid.Parse(item.Uuid)
		foundItem, err := s.inventoryUC.GetItemByID(ctx, uuidItem)
		if err != nil {
			return nil, err
		} else if foundItem != nil {
			foundItem.Qty = foundItem.Qty - 1
			status = pb.Status_OK

			_, err := s.inventoryUC.RemoveItem(ctx, foundItem)
			if err != nil {
				s.logger.Error(err)
				status = pb.Status_Error
				response.StatusMessage = "Error during removing item"
			}
		} else {
			status = pb.Status_NotFound
			response.StatusMessage = "Not found: " + item.Uuid
		}
	}
	return response, nil
}
