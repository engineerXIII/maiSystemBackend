package product

import (
	"context"
	"github.com/engineerXIII/maiSystemBackend/config"
	"github.com/engineerXIII/maiSystemBackend/pkg/logger"
	"github.com/go-co-op/gocron"
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server struct
type Server struct {
	echo        *echo.Echo
	cfg         *config.Config
	amqqChannel *amqp.Channel
	amqpQueue   *amqp.Queue
	scheduler   *gocron.Scheduler
	logger      logger.Logger
}

// NewServer New Server constructor
func NewServer(cfg *config.Config, amqqChannel *amqp.Channel, amqpQueue *amqp.Queue, scheduler *gocron.Scheduler, logger logger.Logger) *Server {
	return &Server{echo: echo.New(), cfg: cfg, amqqChannel: amqqChannel, amqpQueue: amqpQueue, scheduler: scheduler, logger: logger}
}

const (
	certFile       = "ssl/Server.crt"
	keyFile        = "ssl/Server.pem"
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

func (s *Server) Run() error {
	messages, err := s.amqqChannel.Consume(
		s.amqpQueue.Name, // queue
		"",               // consumer
		true,             // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	if err != nil {
		s.logger.Fatalf("AMQP failed to register a consumer. Error: %s", err)
	}

	go func() {
		for message := range messages {
			s.logger.Infof("AMQP received a message: %s", message.Body)
		}
	}()

	if s.cfg.Server.SSL {
		if err := s.MapHandlers(s.echo); err != nil {
			return err
		}

		s.echo.Server.ReadTimeout = time.Second * s.cfg.Server.ReadTimeout
		s.echo.Server.WriteTimeout = time.Second * s.cfg.Server.WriteTimeout

		go func() {
			s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
			s.echo.Server.ReadTimeout = time.Second * s.cfg.Server.ReadTimeout
			s.echo.Server.WriteTimeout = time.Second * s.cfg.Server.WriteTimeout
			s.echo.Server.MaxHeaderBytes = maxHeaderBytes
			if err := s.echo.StartTLS(s.cfg.Server.Port, certFile, keyFile); err != nil {
				s.logger.Fatalf("Error starting TLS Server: ", err)
			}
		}()

		go func() {
			s.logger.Infof("Starting Debug Server on PORT: %s", s.cfg.Server.PprofPort)
			if err := http.ListenAndServe(s.cfg.Server.PprofPort, http.DefaultServeMux); err != nil {
				s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
			}
		}()

		s.scheduler.StartAsync()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
		defer s.scheduler.Stop()
		defer shutdown()

		s.logger.Info("Server Exited Properly")
		return s.echo.Server.Shutdown(ctx)
	}

	server := &http.Server{
		Addr:           s.cfg.Server.Port,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.logger.Fatalf("Error starting Server: ", err)
		}
	}()

	go func() {
		s.logger.Infof("Starting Debug Server on PORT: %s", s.cfg.Server.PprofPort)
		if err := http.ListenAndServe(s.cfg.Server.PprofPort, http.DefaultServeMux); err != nil {
			s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
		}
	}()

	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	s.scheduler.StartAsync()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer s.scheduler.Stop()
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}
