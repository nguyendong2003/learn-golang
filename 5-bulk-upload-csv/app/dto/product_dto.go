package dto

import (
	"bulk-upload-csv/model"
	"time"
)

type ProductResponse struct {
	ID         string    `json:"id"`
	Sku        string    `json:"sku"`
	CategoryID string    `json:"categoryId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func NewListProductResponse(data []model.Product) []ProductResponse {
	// Khởi tạo slice với độ dài cố định để tối ưu bộ nhớ
	result := make([]ProductResponse, len(data))
	for i, product := range data {
		result[i] = NewProductDetailResponse(product)
	}
	return result
}

func NewProductDetailResponse(data model.Product) ProductResponse {
	return ProductResponse{
		ID:         data.ID.String(),
		Sku:        data.Sku,
		CategoryID: data.CategoryID.String(),
		CreatedAt:  data.CreatedAt,
		UpdatedAt:  data.UpdatedAt,
	}
}

type GetListProductRequest struct {
	PagingRequest
}

type GetProductDetailRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type CreateProductRequest struct {
	Sku        string `json:"sku" binding:"required"`
	CategoryID string `json:"categoryId" binding:"required,uuid"`
}
