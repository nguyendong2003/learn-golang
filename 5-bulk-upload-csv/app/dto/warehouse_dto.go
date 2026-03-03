package dto

import (
	"bulk-upload-csv/model"
	"time"
)

type WarehouseResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewListWarehouseResponse(data []model.Warehouse) []WarehouseResponse {
	// Khởi tạo slice với độ dài cố định để tối ưu bộ nhớ
	result := make([]WarehouseResponse, len(data))
	for i, Warehouse := range data {
		result[i] = NewWarehouseDetailResponse(Warehouse)
	}
	return result
}

func NewWarehouseDetailResponse(data model.Warehouse) WarehouseResponse {
	return WarehouseResponse{
		ID:        data.ID.String(),
		Code:      data.Code,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

type GetListWarehouseRequest struct {
	PagingRequest
}

type GetWarehouseDetailRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type CreateWarehouseRequest struct {
	Code string `json:"code" binding:"required"`
}
