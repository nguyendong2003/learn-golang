package dto

import (
	"time"
)

type Pagination struct {
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	TotalItems  int  `json:"total"`
	HasNext     bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
}

type ApiResponse[T any] struct {
	Message    string      `json:"message,omitempty"`
	Data       *T          `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Timestamp  string      `json:"timestamp"`
}

func NewPagination(
	page int,
	limit int,
	totalItems int,
) *Pagination {

	totalPages := 0
	if limit > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}

	return &Pagination{
		Page:        page,
		Limit:       limit,
		TotalItems:  totalItems,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}
}

func newBaseResponse[T any](message string) ApiResponse[T] {
	return ApiResponse[T]{
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func SimpleSuccess(message string) ApiResponse[any] {
	if message == "" {
		message = "Request successful"
	}
	res := newBaseResponse[any](message)
	return res
}

func Success[T any](data *T, message string) ApiResponse[T] {
	if message == "" {
		message = "Request successful"
	}
	res := newBaseResponse[T](message)
	res.Data = data
	return res
}

func Paginated[T any](data *T, p *Pagination, message string) ApiResponse[T] {
	res := Success(data, message)
	res.Pagination = p
	return res
}
