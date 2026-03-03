package dto

import (
	"bulk-upload-csv/model"
	"time"
)

type CategoryResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewListCategoryResponse(data []model.Category) []CategoryResponse {
	// Khởi tạo slice với độ dài cố định để tối ưu bộ nhớ
	result := make([]CategoryResponse, len(data))
	for i, category := range data {
		result[i] = NewCategoryDetailResponse(category)
	}
	return result
}

func NewCategoryDetailResponse(data model.Category) CategoryResponse {
	return CategoryResponse{
		ID:        data.ID.String(),
		Code:      data.Code,
		Name:      data.Name,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

type GetListCategoryRequest struct {
	PagingRequest
}

type GetCategoryDetailRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type CreateCategoryRequest struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}
