package product

import (
	"fmt"
	"github.com/engineerXIII/maiSystemBackend/docs"

	//authRepository "github.com/engineerXIII/maiSystemBackend/internal/auth/repository"
	//authUseCase "github.com/engineerXIII/maiSystemBackend/internal/auth/usecase"
	apiMiddlewares "github.com/engineerXIII/maiSystemBackend/internal/middleware"
	notificationScheduler "github.com/engineerXIII/maiSystemBackend/internal/notification/scheduler"
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

	notificationCron := notificationScheduler.NewNotificationScheduler(s.cfg, s.amqqChannel, s.amqpQueue, s.logger)
	notificationCron.MapCron(s.scheduler)

	mw := apiMiddlewares.NewMiddlewareManager(nil, nil, s.cfg, []string{"*"}, s.logger)

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
	//orderGroup := v1.Group("/order")
	//authGroup := v1.Group("/auth")
	//productGroup := v1.Group("/product")
	//newsGroup := v1.Group("/news")
	//commGroup := v1.Group("/comments")

	//orderHttp.MapOrderRoutes(orderGroup, orderHandlers, mw)
	//authHttp.MapAuthRoutes(authGroup, authHandlers, mw)
	//productHttp.MapProductRoutes(productGroup, productHandlers, mw)
	//newsHttp.MapNewsRoutes(newsGroup, newsHandlers, mw)
	//commentsHttp.MapCommentsRoutes(commGroup, commHandlers, mw)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
