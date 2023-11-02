package product

import (
	"context"
	"github.com/engineerXIII/maiSystemBackend/config"
	"github.com/engineerXIII/maiSystemBackend/pkg/logger"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	grpcServer  *grpc.Server
	redisClient *redis.Client
	logger      logger.Logger
}

// NewServer New Server constructor
func NewServer(cfg *config.Config, redisClient *redis.Client, logger logger.Logger) *Server {
	return &Server{echo: echo.New(), cfg: cfg, redisClient: redisClient, logger: logger}
}

const (
	certFile       = "ssl/cert.crt"
	keyFile        = "ssl/private.pem"
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

func (s *Server) Run() error {
	grpcOpts := []grpc.ServerOption{}
	if true {
		//if s.cfg.Server.SSL {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			return err
		}
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}

	s.grpcServer = grpc.NewServer(grpcOpts...)

	if s.cfg.Server.SSL {
		if err := s.MapHandlers(s.echo); err != nil {
			return err
		}
		s.echo.Server.ReadTimeout = time.Second * s.cfg.Server.ReadTimeout
		s.echo.Server.WriteTimeout = time.Second * s.cfg.Server.WriteTimeout
		s.echo.Server.MaxHeaderBytes = maxHeaderBytes
		go func() {
			s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
			if err := s.echo.StartTLS(s.cfg.Server.Port, certFile, keyFile); err != nil {
				s.logger.Fatalf("Error starting TLS Server: ", err)
			}
		}()

		if s.cfg.Server.Debug {
			go func() {
				s.logger.Infof("Starting Debug Server on PORT: %s", s.cfg.Server.PprofPort)
				if err := http.ListenAndServe(s.cfg.Server.PprofPort, http.DefaultServeMux); err != nil {
					s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
				}
			}()
		}

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		<-quit

		ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
		defer shutdown()
		defer s.grpcServer.GracefulStop()

		s.logger.Info("Server Exited Properly")
		return s.echo.Shutdown(ctx)
	}

	if err := s.MapHandlers(s.echo); err != nil {
		return err
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

	if s.cfg.Server.Debug {
		go func() {
			s.logger.Infof("Starting Debug Server on PORT: %s", s.cfg.Server.PprofPort)
			if err := http.ListenAndServe(s.cfg.Server.PprofPort, http.DefaultServeMux); err != nil {
				s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
			}
		}()
	}
	//
	//if err := s.MapHandlers(s.echo); err != nil {
	//	return err
	//}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()
	defer s.grpcServer.GracefulStop()

	s.logger.Info("Server Exited Properly")
	return s.echo.Shutdown(ctx)
}
