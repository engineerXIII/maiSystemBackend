package product

import (
	"fmt"
	"github.com/engineerXIII/maiSystemBackend/docs"
	authHttp "github.com/engineerXIII/maiSystemBackend/internal/auth/delivery/http"
	authRepository "github.com/engineerXIII/maiSystemBackend/internal/auth/repository"
	authUseCase "github.com/engineerXIII/maiSystemBackend/internal/auth/usecase"
	apiMiddlewares "github.com/engineerXIII/maiSystemBackend/internal/middleware"
	productRepository "github.com/engineerXIII/maiSystemBackend/internal/product/repository"
	productUseCase "github.com/engineerXIII/maiSystemBackend/internal/product/usecase"
	sessionRepository "github.com/engineerXIII/maiSystemBackend/internal/session/repository"
	"github.com/engineerXIII/maiSystemBackend/internal/session/usecase"
	"github.com/engineerXIII/maiSystemBackend/pkg/csrf"
	"github.com/engineerXIII/maiSystemBackend/pkg/metric"
	"github.com/engineerXIII/maiSystemBackend/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
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
	aRepo := authRepository.NewAuthRepository(s.db)
	pRepo := productRepository.NewProductRepository(s.db)
	//nRepo := newsRepository.NewNewsRepository(s.db)
	//cRepo := commentsRepository.NewCommentsRepository(s.db)
	//aAWSRepo := authRepository.NewAuthAWSRepository(s.awsClient)
	authRedisRepo := authRepository.NewAuthRedisRepo(s.redisClient)
	//newsRedisRepo := newsRepository.NewNewsRedisRepo(s.redisClient)

	// Init useCases
	authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, authRedisRepo, s.logger)
	pUC := productUseCase.NewProductUseCase(s.cfg, pRepo, s.logger)
	//newsUC := newsUseCase.NewNewsUseCase(s.cfg, nRepo, newsRedisRepo, s.logger)
	//commUC := commentsUseCase.NewCommentsUseCase(s.cfg, cRepo, s.logger)
	sessUC := usecase.NewSessionUseCase(sRepo, s.cfg)

	// Init handlers
	authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC, sessUC, s.logger)
	//newsHandlers := newsHttp.NewNewsHandlers(s.cfg, newsUC, s.logger)
	//commHandlers := commentsHttp.NewCommentsHandlers(s.cfg, commUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(sessUC, authUC, s.cfg, []string{"*"}, s.logger)

	e.Use(mw.RequestLoggerMiddleware)

	if s.cfg.Docs.Enable {
		docs.SwaggerInfo.Title = s.cfg.Docs.Title
		e.GET(fmt.Sprintf("/%s/*", s.cfg.Docs.Prefix), echoSwagger.WrapHandler)
	}

	if s.cfg.Server.SSL {
		e.Pre(middleware.HTTPSRedirect())
	}

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

	v1 := e.Group("/api/v1")

	health := v1.Group("/health")
	authGroup := v1.Group("/auth")
	//newsGroup := v1.Group("/news")
	//commGroup := v1.Group("/comments")

	authHttp.MapAuthRoutes(authGroup, authHandlers, mw)
	//newsHttp.MapNewsRoutes(newsGroup, newsHandlers, mw)
	//commentsHttp.MapCommentsRoutes(commGroup, commHandlers, mw)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
