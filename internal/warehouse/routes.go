package warehouse

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *WarehouseHandler) {
	router.POST("", handler.CreateWarehouse)
	router.GET("", handler.ListWarehouses)
}
