package scheduler

import (
	"context"
	"github.com/engineerXIII/maiSystemBackend/config"
	"github.com/engineerXIII/maiSystemBackend/internal/notification"
	"github.com/engineerXIII/maiSystemBackend/pkg/logger"
	"github.com/go-co-op/gocron"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type notificationScheduler struct {
	cfg         *config.Config
	amqqChannel *amqp.Channel
	amqpQueue   *amqp.Queue
	logger      logger.Logger
}

func NewNotificationScheduler(cfg *config.Config, amqqChannel *amqp.Channel, amqpQueue *amqp.Queue, logger logger.Logger) order.Scheduler {
	return &notificationScheduler{cfg: cfg, amqqChannel: amqqChannel, amqpQueue: amqpQueue, logger: logger}
}

func (o *notificationScheduler) MapCron(cron *gocron.Scheduler) {
	cron.Every(30).Second().Do(func() {
		ctx, shutdown := context.WithTimeout(context.Background(), 30*time.Second)
		ctx.Done()
		defer shutdown()
		//
		//keys, err := repo.GetOrderKeysCtx(context.Background())
		//if err != nil {
		//	o.logger.Errorf("[CRON][AUTOSTATUS]: Order keys scan redis failed: %s", err)
		//}
		//
		//var keyLen int = len(keys)
		//if keyLen == 0 {
		//	o.logger.Info("[CRON][AUTOSTATUS]: Nothing to update in orders")
		//	return
		//}
		//index_key := rand.Intn(keyLen)
		//value, err := repo.GetOrderByIDCtx(ctx, keys[index_key])
		//switch value.Status {
		//default:
		//	return
		//case models.OrderStatusPackaged:
		//	value.Status = value.Status + 1
		//	value.StatusMessage = value.Status.ToString()
		//	break
		//case models.OrderStatusCreated:
		//	value.Status = value.Status + 1
		//	value.StatusMessage = value.Status.ToString()
		//	break
		//}
		//err = repo.SetOrderCtx(ctx, keys[index_key], 3600, value)
		//if err != nil {
		//	o.logger.Errorf("[CRON][AUTOSTATUS]: Order update fail: %s", err)
		//}
		//
		//jsonStr, _ := json.Marshal(models.OrderStatusNotify{
		//	OrderId:       value.OrderId,
		//	Status:        value.Status,
		//	StatusMessage: value.StatusMessage,
		//})
		//
		//o.logger.Debugf("[CRON][AUTOSTATUS]: Order notify JSON: %s", string(jsonStr))
		//
		//err = o.amqqChannel.PublishWithContext(ctx,
		//	"",
		//	o.amqpQueue.Name,
		//	false,
		//	false,
		//	amqp.Publishing{
		//		ContentType: "application/json",
		//		Body:        jsonStr,
		//	})
		//if err != nil {
		//	o.logger.Errorf("[CRON][AUTOSTATUS]: Order notify publish error: %s", err)
		//}

	})
}
