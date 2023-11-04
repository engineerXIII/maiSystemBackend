package http

import (
	"github.com/engineerXIII/maiSystemBackend/config"
	"github.com/engineerXIII/maiSystemBackend/internal/models"
	"github.com/engineerXIII/maiSystemBackend/internal/order"
	"github.com/engineerXIII/maiSystemBackend/pkg/httpErrors"
	"github.com/engineerXIII/maiSystemBackend/pkg/logger"
	"github.com/engineerXIII/maiSystemBackend/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type orderHandlers struct {
	cfg     *config.Config
	orderUC order.UseCase
	logger  logger.Logger
}

func NewOrderHandlers(cfg *config.Config, orderUC order.UseCase, logger logger.Logger) order.Handlers {
	return &orderHandlers{cfg: cfg, orderUC: orderUC, logger: logger}
}

// Create godoc
// @Summary Create order
// @Description Create order handler
// @Tags Order
// @Accept json
// @Produce json
// @Success 201 {object} models.Order
// @Router /order/create [post]
func (h orderHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "orderHandlers.Create")
		defer span.Finish()

		p := &models.Order{}
		if err := c.Bind(p); err != nil {
			h.logger.Error(err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		order, err := h.orderUC.Create(ctx, p)
		if err != nil {
			h.logger.Error(err)
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, order)
	}
}

// Update godoc
// @Summary Update order
// @Description Update order handler
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "order_id"
// @Success 201 {object} models.Order
// @Router /order/{id} [put]
func (h orderHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "orderHandlers.Update")
		defer span.Finish()

		orderUUID, err := uuid.Parse(c.Param("order_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		p := &models.Order{}
		if err := c.Bind(p); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		p.OrderId = orderUUID

		updatedOrder, err := h.orderUC.Update(ctx, p)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedOrder)
	}
}

// GetByID godoc
// @Summary Get by id order
// @Description Get by id order handler
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "order_id"
// @Success 200 {object} models.Order
// @Router /order/{id} [get]
func (h orderHandlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "orderHandlers.GetByID")
		defer span.Finish()

		orderUUID, err := uuid.Parse(c.Param("order_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		p, err := h.orderUC.GetOrderByID(ctx, orderUUID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, p)
	}
}

// Delete godoc
// @Summary Delete order
// @Description Delete order handler
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "order_id"
// @Success 200 {string} string "ok"
// @Router /order/{id} [delete]
func (h orderHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "orderHandlers.Delete")
		defer span.Finish()

		orderUUID, err := uuid.Parse(c.Param("order_id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err = h.orderUC.Delete(ctx, orderUUID); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}
