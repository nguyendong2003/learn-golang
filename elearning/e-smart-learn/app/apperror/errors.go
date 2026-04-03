package apperror

import (
	"net/http"
)

type ErrorType string

const (
	NotFound        ErrorType = "not_found"
	BadRequest      ErrorType = "bad_request"
	InternalServer  ErrorType = "internal_server"
	Unauthorized    ErrorType = "unauthorized"
	Forbidden       ErrorType = "forbidden"
	DuplicateEntry  ErrorType = "duplicate_entry"
	ValidationError ErrorType = "validation_error"
)

type AppError struct {
	Type    ErrorType         `json:"type"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
	Err     error             `json:"-"`
}

func MapErrorToStatus(t ErrorType) int {
	switch t {
	case BadRequest, ValidationError:
		return http.StatusBadRequest
	case NotFound:
		return http.StatusNotFound
	case Unauthorized:
		return http.StatusUnauthorized
	case Forbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func (e *AppError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Type:    NotFound,
		Message: message,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Type:    BadRequest,
		Message: message,
	}
}

func NewInternalServerError(message string) *AppError {
	return &AppError{
		Type:    InternalServer,
		Message: message,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Type:    Unauthorized,
		Message: message,
	}
}
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Type:    Forbidden,
		Message: message,
	}
}

func NewDuplicateEntryError(message string) *AppError {
	return &AppError{
		Type:    DuplicateEntry,
		Message: message,
	}
}

func NewValidationError(fields map[string]string) *AppError {
	return &AppError{
		Type:    ValidationError,
		Message: "Validation failed",
		Fields:  fields,
	}
}
