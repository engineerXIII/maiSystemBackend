package product

import "github.com/labstack/echo/v4"

type Handlers interface {
	Create() echo.HandlerFunc
	Update() echo.HandlerFunc
	GetByID() echo.HandlerFunc
	Delete() echo.HandlerFunc
	GetProducts() echo.HandlerFunc
	SearchByName() echo.HandlerFunc
}
