package helpers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ValidationErrorResponse formats binding errors without exposing Go internals.
func ValidationErrorResponse(err error) ErrorResponse {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		fields := make(map[string]string, len(validationErrors))
		for _, field := range validationErrors {
			fields[field.Field()] = validationMessage(field)
		}
		return ErrorResponse{Error: "validation failed", Fields: fields}
	}

	var typeError *json.UnmarshalTypeError
	if errors.As(err, &typeError) && typeError.Field != "" {
		return ErrorResponse{
			Error:  "validation failed",
			Fields: map[string]string{typeError.Field: "has an invalid type"},
		}
	}
	return ErrorResponse{Error: "invalid JSON body"}
}

func validationMessage(field validator.FieldError) string {
	switch field.Tag() {
	case "required":
		return "is required"
	case "notblank":
		return "must not be blank"
	case "gt":
		return fmt.Sprintf("must be greater than %s", field.Param())
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", field.Param())
	default:
		return "is invalid"
	}
}
