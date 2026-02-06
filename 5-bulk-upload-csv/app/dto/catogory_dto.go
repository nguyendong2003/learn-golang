package dto

import (
	"bulk-upload-csv/model"
	"time"
)

type CategoryResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewListCategoryResponse(data []model.Category) []CategoryResponse {
	var result []CategoryResponse
	for _, category := range data {
		result = append(result, NewCategoryDetailResponse(category))
	}

	return result
}

func NewCategoryDetailResponse(data model.Category) CategoryResponse {
	return CategoryResponse{
		ID:        data.ID.String(),
		Name:      data.Name,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

type GetListCategoryRequest struct {
	PagingRequest
}

type GetCategoryDetailRequest struct {
	Id string `param:"id"`
}

type CreateCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
