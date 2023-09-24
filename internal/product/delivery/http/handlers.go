package http

import (
	"github.com/engineerXIII/maiSystemBackend/config"
	"github.com/engineerXIII/maiSystemBackend/internal/models"
	"github.com/engineerXIII/maiSystemBackend/internal/product"
	"github.com/engineerXIII/maiSystemBackend/pkg/httpErrors"
	"github.com/engineerXIII/maiSystemBackend/pkg/logger"
	"github.com/engineerXIII/maiSystemBackend/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type productHandlers struct {
	cfg       *config.Config
	productUC product.UseCase
	logger    logger.Logger
}

func NewProductHandlers(cfg *config.Config, productUC product.UseCase, logger logger.Logger) product.Handlers {
	return &productHandlers{cfg: cfg, productUC: productUC, logger: logger}
}

// Create godoc
// @Summary Create product
// @Description Create product handler
// @Tags Product
// @Accept json
// @Produce json
// @Success 201 {object} models.Product
// @Router /product/create [post]
func (h productHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "productHandlers.Create")
		defer span.Finish()

		p := &models.Product{}
		if err := c.Bind(p); err != nil {
			h.logger.Error(err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		createdProduct, err := h.productUC.Create(ctx, p)
		if err != nil {
			h.logger.Error(err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, createdProduct)
	}
}

// Update godoc
// @Summary Update product
// @Description Update product handler
// @Tags Product
// @Accept json
// @Produce json
// @Param id path int true "product_id"
// @Success 201 {object} models.Product
// @Router /product/{id} [put]
func (h productHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "productHandlers.Update")
		defer span.Finish()

		productUUID, err := uuid.Parse(c.Param("product_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		p := &models.Product{}
		if err := c.Bind(p); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		p.ProductID = productUUID

		updatedProduct, err := h.productUC.Update(ctx, p)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedProduct)
	}
}

// GetByID godoc
// @Summary Get by id product
// @Description Get by id product handler
// @Tags Product
// @Accept json
// @Produce json
// @Param id path int true "product_id"
// @Success 200 {object} models.Product
// @Router /product/{id} [get]
func (h productHandlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "productHandlers.GetByID")
		defer span.Finish()

		productUUID, err := uuid.Parse(c.Param("product_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		p, err := h.productUC.GetProductByID(ctx, productUUID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, p)
	}
}

// Delete godoc
// @Summary Delete product
// @Description Delete product handler
// @Tags Product
// @Accept json
// @Produce json
// @Param id path int true "product_id"
// @Success 200 {string} string "ok"
// @Router /product/{id} [delete]
func (h productHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "productHandlers.Delete")
		defer span.Finish()

		productUUID, err := uuid.Parse(c.Param("product_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err = h.productUC.Delete(ctx, productUUID); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// GetProducts godoc
// @Summary Get product list
// @Description Get product list handler
// @Tags Product
// @Accept json
// @Produce json
// @Param page query int flase "page number" Format(page)
// @Param size query int flase "size of page" Format(size)
// @Param orderBy query int flase "filter name" Format(orderBy)
// @Success 200 {object} models.ProductList
// @Router /product [get]
func (h productHandlers) GetProducts() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "productHandlers.GetProducts")
		defer span.Finish()

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		productList, err := h.productUC.GetProducts(ctx, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, productList)
	}
}

// SearchByName godoc
// @Summary Search product by name
// @Description Search product by name handler
// @Tags Product
// @Accept json
// @Produce json
// @Param page query int flase "page number" Format(page)
// @Param size query int flase "size of page" Format(size)
// @Param orderBy query int flase "filter name" Format(orderBy)
// @Success 200 {object} models.ProductList
// @Router /product/search [get]
func (h productHandlers) SearchByName() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "productHandlers.SearchByName")
		defer span.Finish()

		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		productList, err := h.productUC.SearchByName(ctx, c.QueryParam("name"), pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, productList)
	}
}
