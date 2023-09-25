package http

import (
	"github.com/engineerXIII/maiSystemBackend/internal/middleware"
	"github.com/engineerXIII/maiSystemBackend/internal/order"
	"github.com/labstack/echo/v4"
)

func MapOrderRoutes(orderGroup *echo.Group, p order.Handlers, mw *middleware.MiddlewareManager) {
	orderGroup.POST("/create", p.Create()) //, mw.AuthSessionMiddleware, mw.CSRF)
	orderGroup.PUT("/:order_id", p.Update())
	orderGroup.DELETE("/:order_id", p.Delete())
	orderGroup.GET("/:order_id", p.GetByID())
}
