package dto

import (
	"elearning-api/model"
	"time"
)

type CourseProgress struct {
	Percentage       int        `json:"percentage"`
	CompletedLessons int        `json:"completed_lessons"`
	TotalLessons     int        `json:"total_lessons"`
	LastAccessedAt   *time.Time `json:"last_accessed_at,omitempty"`
}

type MyCourseItem struct {
	Course     *CourseResponse `json:"course"`
	Progress   CourseProgress  `json:"progress"`
	EnrolledAt time.Time       `json:"enrolled_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

func NewMyCourseItem(e *model.Enrollment) *MyCourseItem {
	courseDTO := NewCourseDetailResponse(e.Course)

	// compute total lessons by summing lessons in chapters
	total := 0
	for _, ch := range e.Course.Chapters {
		if ch != nil {
			total += len(ch.Lessons)
		}
	}

	completed := len(e.LearnedLessonIds)
	pct := 0
	if total > 0 {
		pct = (completed * 100) / total
	}

	return &MyCourseItem{
		Course: courseDTO,
		Progress: CourseProgress{
			Percentage:       pct,
			CompletedLessons: completed,
			TotalLessons:     total,
		},
		EnrolledAt: e.EnrolledAt,
		UpdatedAt:  e.EnrolledAt,
	}
}

type CourseProgresssRequest struct {
	PagingRequest
}
