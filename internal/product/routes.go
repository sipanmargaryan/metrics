package product

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *ProductHandler) {
	router.POST("", handler.CreateProduct)
	router.GET("", handler.ListProducts)
}
