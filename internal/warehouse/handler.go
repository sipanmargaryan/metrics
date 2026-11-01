package warehouse

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"metrics/internal/helpers"
)

type WarehouseHandler struct {
	service WarehouseService
}

func InitWarehouseHandler(service WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service: service}
}

// CreateWarehouse creates a warehouse.
// @Summary Create a warehouse
// @Tags warehouses
// @Accept json
// @Produce json
// @Param warehouse body CreateWarehouseRequest true "Warehouse to create"
// @Success 201 {object} WarehouseResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 409 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /warehouses [post]
func (handler *WarehouseHandler) CreateWarehouse(c *gin.Context) {
	var request CreateWarehouseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(fmt.Errorf("%w: %w", helpers.ErrValidation, err))
		return
	}
	warehouse, err := handler.service.CreateWarehouse(c.Request.Context(), request)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, toWarehouseResponse(*warehouse))
}

// ListWarehouses lists warehouses.
// @Summary List warehouses
// @Tags warehouses
// @Produce json
// @Success 200 {array} WarehouseResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /warehouses [get]
func (handler *WarehouseHandler) ListWarehouses(c *gin.Context) {
	warehouses, err := handler.service.ListWarehouses(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, helpers.Map(warehouses, toWarehouseResponse))
}
