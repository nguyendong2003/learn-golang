package dto

import (
	"time"

	"elearning-api/model"

	"github.com/google/uuid"
)

type CreateFeedbackRequest struct {
	CourseID uuid.UUID `json:"course_id" binding:"required"`
	Rating   int       `json:"rating" binding:"required,min=1,max=5"`
	Comment  string    `json:"comment"`
}

type FeedbackResponse struct {
	ID        string        `json:"id"`
	CourseID  string        `json:"course_id"`
	Rate      int           `json:"rate"`
	Content   string        `json:"content"`
	User      *UserResponse `json:"user_id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type TestimonialResponse struct {
	ID          uuid.UUID `json:"id"`
	Image       string    `json:"image"`
	Name        string    `json:"name"`
	Comment     string    `json:"comment"`
	Rating      int       `json:"rating"`
	ReviewCount int64     `json:"reviewCount"`
}

func NewFeedbackResponse(data *model.Feedback) *FeedbackResponse {
	if data == nil {
		return nil
	}
	return &FeedbackResponse{
		ID:        data.ID.String(),
		CourseID:  data.CourseID.String(),
		Rate:      data.Rate,
		Content:   data.Content,
		User:      NewUserDetailResponse(data.User),
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func NewFeedbackResponses(data []*model.Feedback) []*FeedbackResponse {
	responses := make([]*FeedbackResponse, len(data))
	for i, feedback := range data {
		responses[i] = NewFeedbackResponse(feedback)
	}
	return responses
}
