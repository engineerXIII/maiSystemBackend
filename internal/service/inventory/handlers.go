package product

import (
	"github.com/engineerXIII/maiSystemBackend/docs"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	inventoryHandler "github.com/engineerXIII/maiSystemBackend/internal/inventory/delivery/grpc"
	inventoryRepository "github.com/engineerXIII/maiSystemBackend/internal/inventory/repository"
	pb "github.com/engineerXIII/maiSystemBackend/proto/api/v1"
	//authRepository "github.com/engineerXIII/maiSystemBackend/internal/auth/repository"
	//authUseCase "github.com/engineerXIII/maiSystemBackend/internal/auth/usecase"
	inventoryUsecase "github.com/engineerXIII/maiSystemBackend/internal/inventory/usecase"
	apiMiddlewares "github.com/engineerXIII/maiSystemBackend/internal/middleware"
	sessionRepository "github.com/engineerXIII/maiSystemBackend/internal/session/repository"
	seccUseCase "github.com/engineerXIII/maiSystemBackend/internal/session/usecase"
	"github.com/engineerXIII/maiSystemBackend/pkg/csrf"
	"github.com/engineerXIII/maiSystemBackend/pkg/metric"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strings"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	metrics, err := metric.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.ServiceName)
	if err != nil {
		s.logger.Errorf("CreateMetrics Error: %s", err)
	}
	s.logger.Infof(
		"Metrics available URL: %s, ServiceName: %s",
		s.cfg.Metrics.URL,
		s.cfg.Metrics.ServiceName,
	)

	// Init repositories
	sRepo := sessionRepository.NewSessionRepository(s.redisClient, s.cfg)
	iRepo := inventoryRepository.NewInventoryRedisRepo(s.redisClient)
	////aRepo := authRepository.NewAuthRepository(s.db)
	//orderRedisRepo := orderRepository.NewOrderRedisRepo(s.redisClient)
	//
	//// Init useCases
	////authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, authRedisRepo, s.logger)
	sessUC := seccUseCase.NewSessionUseCase(sRepo, s.cfg)
	inventoryUC := inventoryUsecase.NewInventoryUseCase(s.cfg, iRepo, s.logger)
	//orderUC := orderUseCase.NewOrderUseCase(s.cfg, orderRedisRepo, s.logger)

	// Init handlers
	//authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC, sessUC, s.logger)
	//orderHandlers := orderHttp.NewOrderHandlers(s.cfg, orderUC, s.logger)
	//
	//orderScheduler := orderScheduler.NewOrderScheduler(s.cfg, s.amqqChannel, s.amqpQueue, &orderRedisRepo, s.logger)
	//orderScheduler.MapCron(s.scheduler)
	inventoryHandler := inventoryHandler.NewInventoryServer(s.cfg, inventoryUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(sessUC, nil, s.cfg, []string{"*"}, s.logger)

	e.Use(mw.RequestLoggerMiddleware)

	if s.cfg.Docs.Enable {
		docs.SwaggerInfo.Title = s.cfg.Docs.Title
		//e.GET(fmt.Sprintf("/%s/*", s.cfg.Docs.Prefix), echoSwagger.WrapHandler)
	}

	//if s.cfg.Server.SSL {
	//	e.Pre(middleware.HTTPSRedirect())
	//}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, csrf.CSRFHeader},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())
	e.Use(mw.MetricsMiddleware(metrics))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("3M"))
	if s.cfg.Server.Debug {
		e.Use(mw.DebugMiddleware)
	}

	pb.RegisterInventoryServiceServer(s.grpcServer, inventoryHandler)

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		writer.Write([]byte("{\"status\":\"OK\"}"))
	})
	e.Any("/*", func(c echo.Context) error {
		h2c.NewHandler(
			http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.ProtoMajor == 2 && strings.Contains(req.Header.Get("Content-Type"), "application/grpc") {
					s.logger.Debug("Called grpc")
					s.grpcServer.ServeHTTP(w, req)
				} else {
					s.logger.Debug("Called http")
					httpMux.ServeHTTP(w, req)
				}
			}),
			&http2.Server{}).ServeHTTP(c.Response(), c.Request())
		return nil
	})

	return nil
}
