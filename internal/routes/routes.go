package routes

import (
	"github.com/gin-gonic/gin"

	"metrics/internal/config"
	"metrics/internal/order"
	"metrics/internal/product"
	"metrics/internal/warehouse"
)

func Register(router *gin.Engine, app *config.App) {
	api := router.Group("/api/v1")
	warehouseRepo := warehouse.InitWarehouseRepository(app.DB)
	warehouseService := warehouse.InitWarehouseService(warehouseRepo, app.Config)
	warehouseHandler := warehouse.InitWarehouseHandler(warehouseService)
	warehouse.RegisterRoutes(api.Group("/warehouses"), warehouseHandler)

	productRepo := product.InitProductRepository(app.DB)
	productService := product.InitProductService(productRepo, app.Config)
	productHandler := product.InitProductHandler(productService)
	product.RegisterRoutes(api.Group("/products"), productHandler)

	orderRepo := order.InitOrderRepository(app.DB)
	orderService := order.InitOrderService(orderRepo, app.Config)
	orderHandler := order.InitOrderHandler(orderService)

	order.RegisterRoutes(
		api.Group("/orders"),
		orderHandler,
	)
}
