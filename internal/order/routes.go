package order

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *OrderHandler) {
	router.POST("", handler.CreateOrder)
}
