package job

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

func NewEventNotificationTask(eventID, courseID uuid.UUID, processAt time.Time) (*asynq.Task, []asynq.Option) {
	payload, _ := json.Marshal(EventNotificationPayload{
		EventID:  eventID,
		CourseID: courseID,
		NotifyAt: processAt.UTC(),
	})

	return asynq.NewTask(TypeEventNotification, payload), []asynq.Option{
		asynq.ProcessAt(processAt),
		asynq.TaskID(fmt.Sprintf("notify:%s:%d", eventID, processAt.UTC().Unix())),
	}
}

func NewBlogPublishTask(blogID uuid.UUID, publishAt time.Time) (*asynq.Task, []asynq.Option) {
	payload, _ := json.Marshal(BlogPublishPayload{
		BlogID:    blogID,
		PublishAt: publishAt.UTC(),
	})

	return asynq.NewTask(TypeBlogPublish, payload), []asynq.Option{
		asynq.ProcessAt(publishAt),
		asynq.TaskID(fmt.Sprintf("blog_publish:%s:%d", blogID, publishAt.UTC().Unix())),
	}
}
