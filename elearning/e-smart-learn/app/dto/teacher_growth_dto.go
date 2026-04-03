package dto

type TeacherGrowthTopCategoryResponse struct {
	Name     string  `json:"name"`
	Count    int64   `json:"count"`
	SharePct float64 `json:"share_pct"`
}

type TeacherGrowthStatisticsResponse struct {
	TotalVerifiedTeachers int64                             `json:"total_verified_teachers"`
	NewThisQuarter        int64                             `json:"new_this_quarter"`
	TopCategory           *TeacherGrowthTopCategoryResponse `json:"top_category,omitempty"`
}
