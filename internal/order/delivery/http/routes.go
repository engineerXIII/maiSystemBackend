package http

import (
	"github.com/engineerXIII/maiSystemBackend/internal/middleware"
	"github.com/engineerXIII/maiSystemBackend/internal/product"
	"github.com/labstack/echo/v4"
)

func MapProductRoutes(productGroup *echo.Group, p product.Handlers, mw *middleware.MiddlewareManager) {
	productGroup.POST("/create", p.Create()) //, mw.AuthSessionMiddleware, mw.CSRF)
	productGroup.PUT("/:product_id", p.Update(), mw.AuthSessionMiddleware, mw.CSRF)
	productGroup.DELETE("/:product_id", p.Delete(), mw.AuthSessionMiddleware, mw.CSRF)
	productGroup.GET("/:product_id", p.GetByID())
	productGroup.GET("/search", p.SearchByName())
	productGroup.GET("", p.GetProducts())
}
