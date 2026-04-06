package dto

import "elearning-api/consts"

type CreatePlanRequest struct {
	Name           string              `json:"name" binding:"required,min=2,max=100"`
	Description    *string             `json:"description"`
	BillingCycle   consts.BillingCycle `json:"billing_cycle" binding:"required,oneof=monthly yearly"`
	Price          float64             `json:"price" binding:"required"`
	Currency       *string             `json:"currency" binding:"omitempty,min=3,max=10"`
	Tag            *string             `json:"tag"`
	IsRecommend    *bool               `json:"is_recommend"`
	AccessFeatures []string            `json:"access_features" binding:"omitempty"`
	IsActive       *bool               `json:"is_active"`
}

type UpdatePlanRequest struct {
	Name           *string              `json:"name" binding:"omitempty,min=2,max=100"`
	Description    *string              `json:"description"`
	BillingCycle   *consts.BillingCycle `json:"billing_cycle" binding:"omitempty,oneof=monthly yearly"`
	Price          *float64             `json:"price"`
	Currency       *string              `json:"currency" binding:"omitempty,min=3,max=10"`
	Tag            *string              `json:"tag"`
	IsRecommend    *bool                `json:"is_recommend"`
	AccessFeatures []string             `json:"access_features" binding:"omitempty"`
}

type ListPlanRequest struct {
	PagingRequest

	IsActive    *bool `form:"is_active,omitempty"`
	IsRecommend *bool `form:"is_recommend,omitempty"`
}
