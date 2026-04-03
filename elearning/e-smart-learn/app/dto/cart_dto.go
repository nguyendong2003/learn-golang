package dto

import "elearning-api/model"

type CartCheckoutRequest struct {
	CouponCode string `json:"coupon_code" binding:"omitempty,max=100"`
}

// type CartItemResponse struct {
// 	Course *CourseResponse `json:"course"`
// }

type CartResponse struct {
	Items       []*CourseResponse `json:"items"`
	TotalAmount float64           `json:"total_amount"`
	Currency    string            `json:"currency"`
}

func NewCartResponse(items []*model.CartItem) *CartResponse {
	res := &CartResponse{Items: make([]*CourseResponse, 0, len(items))}

	for _, item := range items {
		if item == nil || item.Course == nil {
			continue
		}
		courseResp := NewCourseDetailResponse(item.Course)
		res.Items = append(res.Items, courseResp)
		res.TotalAmount += item.Course.Price
		if item.Course.StripeCurrency != "" {
			res.Currency = item.Course.StripeCurrency
		}
	}

	if res.Currency == "" {
		res.Currency = "usd"
	}

	return res
}
