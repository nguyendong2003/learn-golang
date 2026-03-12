package dto

import "elearning-api/model"

type PlanResponse struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	AccessFeatures []string `json:"access_features"`
	PriceMonthly   float64  `json:"price_monthly"`
	PriceYearly    float64  `json:"price_yearly"`
	IsDefault      bool     `json:"is_default"`
}

func NewListPlanResponse(plans []*model.Plan) []*PlanResponse {
	res := make([]*PlanResponse, len(plans))
	for i, p := range plans {
		res[i] = NewPlanDetailResponse(p)
	}
	return res
}

func NewPlanDetailResponse(p *model.Plan) *PlanResponse {
	return &PlanResponse{
		ID:             p.ID.String(),
		Name:           p.Name,
		Description:    p.Description,
		AccessFeatures: []string{"Feature A", "Feature B"},
		PriceMonthly:   p.MonthlyPrice,
		PriceYearly:    p.YearlyPrice,
		IsDefault:      p.IsDefault,
	}
}
