package dto

import (
	"elearning-api/consts"
	"elearning-api/model"
	"time"
)

type CourseResponse struct {
	ID           string                     `json:"id"`
	Title        string                     `json:"title"`
	Description  string                     `json:"description"`
	Image        string                     `json:"image"`
	Slug         string                     `json:"slug"`
	Duration     int                        `json:"duration"`
	Price        float64                    `json:"price"`
	OldPrice     float64                    `json:"old_price"`
	Status       consts.CourseStatus        `json:"status"`
	AverageRate  float64                    `json:"average_rate"`
	TotalStudent int64                      `json:"total_student"`
	Category     *CategoryResponse          `json:"category"`
	Instructor   *InstructorProfileResponse `json:"instructor"`
	CreatedAt    time.Time                  `json:"created_at"`
	UpdatedAt    time.Time                  `json:"updated_at"`
}

func NewListCourseResponse(courses []*model.Course) []*CourseResponse {
	res := make([]*CourseResponse, len(courses))
	for i, c := range courses {
		res[i] = NewCourseDetailResponse(c)
	}
	return res
}

func NewCourseDetailResponse(m *model.Course) *CourseResponse {
	if m == nil {
		return nil
	}

	var category *CategoryResponse
	if m.Category != nil {
		category = NewCategoryDetailResponse(m.Category)
	}

	var instructor *InstructorProfileResponse
	if m.InstructorProfile != nil {
		instructor = NewInstructorProfileDetailResponse(m.InstructorProfile)
	}

	return &CourseResponse{
		ID:           m.ID.String(),
		Title:        m.Title,
		Description:  m.Description,
		Image:        m.Image,
		Slug:         m.Slug,
		Duration:     m.Duration,
		Price:        m.Price,
		OldPrice:     m.OldPrice,
		Status:       m.Status,
		AverageRate:  m.AverageRate,
		TotalStudent: m.TotalStudent,
		Category:     category,
		Instructor:   instructor,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

type CreateCourseRequest struct {
	Title       string  `json:"title" binding:"required,min=3,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	CategoryID  string  `json:"category_id" binding:"required,uuid"`
}

type UpdateCourseRequest struct {
	Title       *string  `json:"title" binding:"omitempty,min=3,max=255"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price" binding:"omitempty,gt=0"`
	CategoryID  *string  `json:"category_id" binding:"omitempty,uuid"`
}

type ListCourseRequest struct {
	PagingRequest

	Title      *string              `form:"title"`
	CategoryID *string              `form:"category_id" binding:"omitempty,uuid"`
	Status     *consts.CourseStatus `form:"status"`
}

type CourseSlugRequest struct {
	Slug string `uri:"slug" binding:"required"`
}
