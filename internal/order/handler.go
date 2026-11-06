package order

import (
	"fmt"
	"metrics/internal/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service OrderService
}

func InitOrderHandler(
	service OrderService,
) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

// CreateOrder creates an order.
// @Summary Create an order
// @Tags orders
// @Accept json
// @Produce json
// @Param order body CreateOrderRequest true "Order to create"
// @Success 201 {object} OrderCreatedResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 404 {object} helpers.ErrorResponse
// @Failure 409 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /orders [post]
func (handler *OrderHandler) CreateOrder(c *gin.Context) {
	var request CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(fmt.Errorf("%w: %w", helpers.ErrValidation, err))
		return
	}

	order, err := handler.service.CreateOrder(c.Request.Context(), request)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, OrderCreatedResponse{
		ID: order.ID,
	})

}
