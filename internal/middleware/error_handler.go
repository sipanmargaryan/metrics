package middleware

import (
	"errors"
	"net/http"

	"metrics/internal/helpers"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		switch {
		case errors.Is(err, helpers.ErrNotFound):
			c.JSON(http.StatusNotFound, helpers.ErrorResponse{
				Error: err.Error(),
			})

		case errors.Is(err, helpers.ErrUnauthorized):
			c.JSON(http.StatusUnauthorized, helpers.ErrorResponse{
				Error: err.Error(),
			})

		case errors.Is(err, helpers.ErrForbidden):
			c.JSON(http.StatusForbidden, helpers.ErrorResponse{
				Error: err.Error(),
			})

		case errors.Is(err, helpers.ErrValidation):
			c.JSON(http.StatusBadRequest, helpers.ValidationErrorResponse(err))

		case errors.Is(err, helpers.ErrConflict), errors.Is(err, helpers.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, helpers.ErrorResponse{
				Error: err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, helpers.ErrorResponse{
				Error: "internal server error",
			})
		}
	}
}
