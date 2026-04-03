package worker

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"elearning-api/consts"
	"elearning-api/job"
	"elearning-api/repository"

	"github.com/hibiken/asynq"
)

type BlogScheduledWorkerHandler struct {
	BlogRepository repository.BlogRepository
}

func NewBlogScheduledWorkerHandler(blogRepository repository.BlogRepository) *BlogScheduledWorkerHandler {
	return &BlogScheduledWorkerHandler{
		BlogRepository: blogRepository,
	}
}

func (h *BlogScheduledWorkerHandler) HandlePublish(ctx context.Context, t *asynq.Task) error {
	var p job.BlogPublishPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	blog, err := h.BlogRepository.FindByID(ctx, p.BlogID, nil)
	if err != nil {
		return err
	}
	if blog == nil {
		return nil
	}
	if blog.Status == consts.BlogStatusPublished || blog.Status == consts.BlogStatusDraft || blog.ScheduledAt == nil {
		return nil
	}
	if math.Abs(p.PublishAt.Sub(*blog.ScheduledAt).Seconds()) > 5 {
		return nil
	}
	blog.Status = consts.BlogStatusPublished
	now := time.Now().UTC()
	blog.PublishedAt = &now
	if _, err := h.BlogRepository.Updates(ctx, blog); err != nil {
		return err
	}
	return nil
}
