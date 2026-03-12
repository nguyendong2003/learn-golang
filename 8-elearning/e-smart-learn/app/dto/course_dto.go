package dto

import (
	"elearning-api/model"
)

type CourseResponse struct {
	ID          string              `json:"id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	ImageURL    string              `json:"image_url"`
	Instructor  *InstructorResponse `json:"instructor"`
	AverageRate float64             `json:"average_rate"`
	OldPrice    float64             `json:"old_price"`
	Price       float64             `json:"price"`
	Category    *CategoryResponse   `json:"category"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
}

func NewListCourseResponse(courses []*model.Course) []*CourseResponse {
	res := make([]*CourseResponse, len(courses))
	for i, c := range courses {
		res[i] = NewCourseDetailResponse(c)
	}
	return res
}

func NewCourseDetailResponse(m *model.Course) *CourseResponse {
	var cat *CategoryResponse
	if m.Category != nil {
		cat = NewCategoryDetailResponse(m.Category)
	}

	var ins *InstructorResponse
	if m.Instructor != nil && m.Instructor.User != nil {
		ins = NewInstructorResponse(m.Instructor)
	}

	return &CourseResponse{
		ID:          m.Slug,
		Title:       m.Title,
		Description: m.Description,
		ImageURL:    m.Image,
		Instructor:  ins,
		AverageRate: m.AverageRate,
		OldPrice:    m.OldPrice,
		Price:       m.Price,
		Category:    cat,
		CreatedAt:   m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   m.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
