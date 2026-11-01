package helpers

import "errors"

var (
	ErrConflict           = errors.New("conflict")
	ErrNotFound           = errors.New("not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrValidation         = errors.New("validation failed")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

// ErrorResponse is the shared JSON body for API errors.
type ErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}
