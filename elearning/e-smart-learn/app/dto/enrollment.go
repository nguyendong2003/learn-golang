package dto

import (
	"elearning-api/model"
	"elearning-api/repository/dbtypes"
	"time"
)

type EnrollmentResponse struct {
	ID               string     `json:"id"`
	CourseId         string     `json:"courseId"`
	UserId           string     `json:"userId"`
	EnrollAt         time.Time  `json:"enrollAt"`
	IsCompleted      bool       `json:"isCompleted"`
	CompletedAt      *time.Time `json:"completedAt,omitempty"`
	CanceledAt       *time.Time `json:"canceledAt,omitempty"`
	Type             string     `json:"type"`
	LearnedLessonIds []string   `json:"learnedLessonIds"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func NewEnrollmentResponse(model *model.Enrollment) *EnrollmentResponse {
	if model == nil {
		return nil
	}

	return &EnrollmentResponse{
		ID:               model.ID.String(),
		CourseId:         model.CourseID.String(),
		UserId:           model.UserID.String(),
		EnrollAt:         model.EnrolledAt,
		IsCompleted:      model.IsCompleted,
		CompletedAt:      model.CompletedAt,
		CanceledAt:       model.CanceledAt,
		Type:             model.Type,
		LearnedLessonIds: dbtypes.ToStringSlice(model.LearnedLessonIds),
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
	}
}

type StudentCourseProgressRequest struct {
	CourseId string `json:"courseId" binding:"required,uuid"`
	LessonId string `json:"lessonId" binding:"required,uuid"`
}
