package scheduler

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"github.com/engineerXIII/maiSystemBackend/config"
	"github.com/engineerXIII/maiSystemBackend/internal/models"
	"github.com/engineerXIII/maiSystemBackend/internal/order"
	"github.com/engineerXIII/maiSystemBackend/pkg/logger"
	pb "github.com/engineerXIII/maiSystemBackend/proto/api/v1"
	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"math/rand"
	"os"
	"time"
)

type orderScheduler struct {
	cfg         *config.Config
	orderRepo   *order.RedisRepository
	grpcClient  pb.InventoryServiceClient
	amqqChannel *amqp.Channel
	amqpQueue   *amqp.Queue
	logger      logger.Logger
}

func NewOrderScheduler(cfg *config.Config, amqqChannel *amqp.Channel, amqpQueue *amqp.Queue, orderRepo *order.RedisRepository, logger logger.Logger) order.Scheduler {
	pemServerCA, err := os.ReadFile("ssl/root.pem")
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		logger.Errorf("failed to add server CA's certificate")
	}
	// Create the credentials and return it
	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}
	connCred := credentials.NewTLS(tlsConfig)
	conn, err := grpc.Dial(cfg.Service.Inventory, grpc.WithTransportCredentials(connCred))
	if err != nil {
		logger.Fatalf("GRPC not connect: %v", err)
	}
	client := pb.NewInventoryServiceClient(conn)
	return &orderScheduler{cfg: cfg, grpcClient: client, amqqChannel: amqqChannel, amqpQueue: amqpQueue, orderRepo: orderRepo, logger: logger}
}

func (o *orderScheduler) MapCron(cron *gocron.Scheduler) {
	// Auto status change
	cron.Every(5).Second().Do(func() {
		repo := *o.orderRepo
		ctx, shutdown := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdown()

		keys, err := repo.GetOrderKeysCtx(context.Background())
		if err != nil {
			o.logger.Errorf("[CRON][AUTOSTATUS]: Order keys scan redis failed: %s", err)
		}

		var keyLen int = len(keys)
		if keyLen == 0 {
			o.logger.Debug("[CRON][AUTOSTATUS]: Nothing to update in orders")
			return
		}
		index_key := rand.Intn(keyLen)
		value, err := repo.GetOrderByIDCtx(ctx, keys[index_key])

		switch value.Status {
		default:
			return
		case models.OrderStatusCreated:
			value.Status = value.Status + 1
			value.StatusMessage = value.Status.ToString()
			break
		case models.OrderStatusConfirmed:
			c, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			grpcPaylod := &pb.ItemRequest{Item: []*pb.Item{}}
			for _, item := range value.OrderList {
				grpcPaylod.Item = append(grpcPaylod.Item, &pb.Item{
					Uuid: item.ItemId.String(),
					Qty:  uint64(item.Qty),
				})
			}
			resp, err := o.grpcClient.CheckItem(c, grpcPaylod)
			if err != nil {
				o.logger.Error(err)
				panic(err)
			}
			o.logger.Debugf("Order %v, %v", value.OrderId, resp.Status)
			if resp.Status == pb.Status_OK {
				o.logger.Infof("Order %v fully packaged", value.OrderId)
				value.Status = value.Status + 1
				value.StatusMessage = value.Status.ToString()
			} else {
				for i, item := range value.OrderList {
					uid, _ := uuid.Parse(resp.Items[i].Item.Uuid)
					if item.ItemId == uid {
						if resp.Items[i].Item.Qty == 0 {
							o.logger.Infof("Inventory not have items for order %v. Cancelling...", value.OrderId)
							value.Status = models.OrderStatusCancelled
							break
						}
						value.OrderList[i].Qty = int(resp.Items[i].Item.Qty)
					}
				}
				if value.Status == models.OrderStatusCancelled {
					break
				}
			}

			o.logger.Infof("Inventory give items for order %v", value.OrderId)
			req := &pb.ItemRequest{Item: []*pb.Item{}}
			for _, item := range resp.Items {
				req.Item = append(req.Item, item.Item)
			}
			_, err = o.grpcClient.RemoveItem(ctx, req)
			if err != nil {
				o.logger.Error(err)
			}
			o.logger.Debugf("Order %s, inventory %s", value.OrderId, resp.Status)
			break
		case models.OrderStatusPackaged:
			value.Status = value.Status + 1
			value.StatusMessage = value.Status.ToString()
			break
		case models.OrderStatusInDelivery:
			value.Status = value.Status + 1
			value.StatusMessage = value.Status.ToString()
			break
		case models.OrderStatusCompleted:
			o.logger.Debugf("Order %s removed from processing as completed", value.OrderId)
			repo.DeleteOrderCtx(ctx, keys[index_key])
			return
			//case models.OrderStatusCancelled:
			//	o.logger.Debugf("Order %s removed from processing as cancelled", value.OrderId)
			//	repo.DeleteOrderCtx(ctx, keys[index_key])
			//	return
		}
		err = repo.SetOrderCtx(ctx, keys[index_key], 3600, value)
		if err != nil {
			o.logger.Errorf("[CRON][AUTOSTATUS]: Order update fail: %s", err)
		}

		jsonStr, _ := json.Marshal(models.OrderStatusNotify{
			OrderId:       value.OrderId,
			Status:        value.Status,
			StatusMessage: value.StatusMessage,
		})

		o.logger.Debugf("[CRON][AUTOSTATUS]: Order notify JSON: %s", string(jsonStr))

		err = o.amqqChannel.PublishWithContext(ctx,
			"",
			o.amqpQueue.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        jsonStr,
			})
		if err != nil {
			o.logger.Errorf("[CRON][AUTOSTATUS]: Order notify publish error: %s", err)
		}

	})
}
