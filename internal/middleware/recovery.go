package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"metrics/internal/helpers"
	"metrics/internal/logging"
)

// RecoverPanics avoids Gin's default request dump, which can contain credentials.
func RecoverPanics() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logging.FromContext(c.Request.Context()).Error("request panic",
					"panic", recovered, "stack", string(debug.Stack()))
				if c.Writer.Written() {
					c.Abort()
					return
				}
				c.AbortWithStatusJSON(http.StatusInternalServerError, helpers.ErrorResponse{
					Error: "internal server error",
				})
			}
		}()
		c.Next()
	}
}
