package product

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"metrics/internal/helpers"
)

type ProductHandler struct {
	service ProductService
}

func InitProductHandler(service ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// CreateProduct creates a product belonging to a warehouse.
// @Summary Create a product
// @Tags products
// @Accept json
// @Produce json
// @Param product body CreateProductRequest true "Product to create"
// @Success 201 {object} ProductResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 404 {object} helpers.ErrorResponse
// @Failure 409 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /products [post]
func (handler *ProductHandler) CreateProduct(c *gin.Context) {
	var request CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(fmt.Errorf("%w: %w", helpers.ErrValidation, err))
		return
	}
	product, err := handler.service.CreateProduct(c.Request.Context(), request)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, toProductResponse(*product))
}

// ListProducts lists products and their warehouse IDs.
// @Summary List products
// @Tags products
// @Produce json
// @Success 200 {array} ProductResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /products [get]
func (handler *ProductHandler) ListProducts(c *gin.Context) {
	products, err := handler.service.ListProducts(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, helpers.Map(products, toProductResponse))
}
