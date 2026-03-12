package handler

import (
	"elearning-api/dto"
	"elearning-api/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// FeedbackHandler provides mocked endpoints for feedbacks.
type FeedbackHandler interface {
	GetFeedbacks() gin.HandlerFunc
}

type feedbackHandler struct{}

// NewFeedbackHandler creates a new FeedbackHandler.
func NewFeedbackHandler() FeedbackHandler { return &feedbackHandler{} }

// GetFeedbacks godoc
// @Summary Get list of feedbacks
// @Description Retrieve list of feedbacks (mocked)
// @Tags feedbacks
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Router /api/v1/feedbacks [get]
func (h *feedbackHandler) GetFeedbacks() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now().UTC()

		f1 := &model.Feedback{
			Rate:    5,
			Content: "Excellent course, very practical and well explained.",
		}
		f1.CreatedAt = now.Add(-72 * time.Hour)
		f1.UpdatedAt = now.Add(-48 * time.Hour)
		f1.User = &model.User{Name: "Student A"}

		f2 := &model.Feedback{
			Rate:    4,
			Content: "Good content but could use more examples.",
		}
		f2.CreatedAt = now.Add(-36 * time.Hour)
		f2.UpdatedAt = now.Add(-12 * time.Hour)
		f2.User = &model.User{Name: "Student B"}

		list := []*model.Feedback{f1, f2}

		resp := dto.NewApiResponse(c)
		resp.Path = "/api/v1/feedbacks"
		resp.Request = dto.GetRequestClient(c)
		resp.Data = list
		resp.Metadata = nil

		c.JSON(http.StatusOK, resp)
	}
}
