package helpers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterValidators must be called at startup before any request validation.
func RegisterValidators() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("unexpected validator engine")
	}
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		return strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
	})
	return v.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
		return fl.Field().Kind() == reflect.String && strings.TrimSpace(fl.Field().String()) != ""
	})
}
