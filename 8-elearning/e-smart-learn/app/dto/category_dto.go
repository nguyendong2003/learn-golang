package dto

import (
	"elearning-api/model"
	"time"
)

type CategoryResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewListCategoryResponse(categories []*model.Category) []*CategoryResponse {
	res := make([]*CategoryResponse, len(categories))
	for i, c := range categories {
		res[i] = NewCategoryDetailResponse(c)
	}
	return res
}

func NewCategoryDetailResponse(data *model.Category) *CategoryResponse {
	if data == nil {
		return nil
	}

	return &CategoryResponse{
		ID:          data.ID.String(),
		Name:        data.Name,
		Description: data.Description,
	}
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=255"`
	Description string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=255"`
	Description string `json:"description"`
}

type ListCategoryRequest struct {
	PagingRequest

	Name *string `form:"name,omitempty"`
}
