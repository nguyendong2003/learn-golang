package dto

import "time"

type RevenueBreakdownResponse struct {
	TotalAmount     float64 `json:"total_amount" gorm:"column:total_amount"`
	InstructorGross float64 `json:"instructor_gross" gorm:"column:instructor_gross"`
	PlatformGross   float64 `json:"platform_gross" gorm:"column:platform_gross"`
	StripeFee       float64 `json:"stripe_fee" gorm:"column:stripe_fee"`
	InstructorNet   float64 `json:"instructor_net" gorm:"column:instructor_net"`
	PlatformNet     float64 `json:"platform_net" gorm:"column:platform_net"`
}

type RevenueStatisticsResponse struct {
	CoursePurchase RevenueBreakdownResponse `json:"course_purchase"`
	Subscription   RevenueBreakdownResponse `json:"subscription"`
	Total          RevenueBreakdownResponse `json:"total"`
}

type RevenueByDayQuery struct {
	Date string `form:"date" binding:"required,datetime=2006-01-02"`
}

type RevenueByMonthQuery struct {
	Year  int `form:"year" binding:"required,min=1970"`
	Month int `form:"month" binding:"required,min=1,max=12"`
}

type RevenueByYearQuery struct {
	Year int `form:"year" binding:"required,min=1970"`
}

type RevenueOverviewResponse struct {
	TotalAmount    float64 `json:"total_amount"`
	PreviousAmount float64 `json:"previous_amount"`
	GrowthPct      float64 `json:"growth_pct"`
}

type RevenueTransactionItemResponse struct {
	TransactionID string              `json:"transaction_id"`
	User          RevenueUserResponse `json:"user"`
	Method        string              `json:"method"`
	Amount        float64             `json:"amount"`
	Type          string              `json:"type"`
	CreatedAt     time.Time           `json:"created_at"`
}

type RevenueUserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
}

type SalesSegmentationItemResponse struct {
	Percent int     `json:"percent"`
	Amount  float64 `json:"amount"`
}

type SalesSegmentationResponse struct {
	MembershipSubs  SalesSegmentationItemResponse `json:"membership_subs"`
	SinglePurchases SalesSegmentationItemResponse `json:"single_purchases"`
}

// TeacherRevenueFilterRequest holds optional date range + pagination query params
// for the GET /api/v1/admin/revenue/statistics/teachers/revenue endpoint.
type TeacherRevenueFilterRequest struct {
	StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string `form:"end_date"   binding:"omitempty,datetime=2006-01-02"`
	Limit     int    `form:"limit"      binding:"omitempty,min=1,max=100"`
	Offset    int    `form:"offset"     binding:"omitempty,min=0"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

func (r *TeacherRevenueFilterRequest) Process() {
	if r.Limit <= 0 {
		r.Limit = 10
	}
	if r.Offset < 0 {
		r.Offset = 0
	}
	if r.SortOrder == "" {
		r.SortOrder = "desc"
	}
}

// TeacherRevenueItemResponse is one row in the teacher revenue list.
// Monetary values are in USD (dollars), converted from cents.
type TeacherRevenueItemResponse struct {
	TeacherID     string  `json:"teacher_id"`
	TeacherName   string  `json:"teacher_name"`
	TeacherEmail  string  `json:"teacher_email"`
	TeacherAvatar string  `json:"teacher_avatar"`
	TotalAmount   float64 `json:"total_amount"`   // gross revenue from course sales
	StripeFee     float64 `json:"stripe_fee"`     // Stripe processing fee share
	InstructorNet float64 `json:"instructor_net"` // total_amount - stripe_fee
	TotalCourses  int64   `json:"total_courses"`  // number of courses owned
}
