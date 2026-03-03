package dto

import (
	"time"
)

type Pagination struct {
	Page         int   `json:"page"`
	PageSize     int   `json:"page_size"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
}

type ErrorDTO struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type Meta struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"request_id,omitempty"`
}

type ApiResponse[T any] struct {
	Message    string      `json:"message,omitempty"`
	Data       *T          `json:"data,omitempty"`
	Error      *ErrorDTO   `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Meta       Meta        `json:"meta"`
}

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func NewSuccessResponse[T any](
	data *T,
	message string,
	requestID string,
) *ApiResponse[T] {
	if message == "" {
		message = "Request successfully"
	}

	return &ApiResponse[T]{
		Message: message,
		Data:    data,
		Meta: Meta{
			Timestamp: now(),
			RequestID: requestID,
		},
	}
}

func NewPaginatedResponse[T any](
	data *T,
	pagination *Pagination,
	message string,
	requestID string,
) *ApiResponse[T] {

	if message == "" {
		message = "Request successfully"
	}

	return &ApiResponse[T]{
		Message:    message,
		Data:       data,
		Pagination: pagination,
		Meta: Meta{
			Timestamp: now(),
			RequestID: requestID,
		},
	}
}

// Code là một chuỗi định danh lỗi, ví dụ: "validation_error", "not_found", "internal_server_error"
func NewErrorResponse(
	code string,
	message string,
	requestID string,
) *ApiResponse[any] {

	if message == "" {
		message = "Something went wrong"
	}

	return &ApiResponse[any]{
		Error: &ErrorDTO{
			Code:    code,
			Message: message,
		},
		Meta: Meta{
			Timestamp: now(),
			RequestID: requestID,
		},
	}
}

func NewValidationErrorResponse(
	fields map[string]string,
	requestID string,
) *ApiResponse[any] {

	return &ApiResponse[any]{
		Error: &ErrorDTO{
			Code:    "VALIDATION_ERROR",
			Message: "Validation failed",
			Fields:  fields,
		},
		Meta: Meta{
			Timestamp: now(),
			RequestID: requestID,
		},
	}
}
