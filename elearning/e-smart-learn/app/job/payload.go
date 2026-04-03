package job

import (
	"time"

	"github.com/google/uuid"
)

const (
	TypeEventNotification = "event:notify"
	TypeSendSingleEmail   = "email:notify_event_reminder"
	TypeBlogPublish       = "blog:publish"
)

type EventNotificationPayload struct {
	EventID  uuid.UUID `json:"event_id"`
	CourseID uuid.UUID `json:"course_id"`
	NotifyAt time.Time `json:"notify_at"`
}

type EmailEventReminderPayload struct {
	Email       string    `json:"email"`
	EventName   string    `json:"event_name"`
	StudentName string    `json:"student_name"`
	StartTime   time.Time `json:"start_time"`
	RoomToken   string    `json:"room_token"`
}

type BlogPublishPayload struct {
	BlogID    uuid.UUID `json:"blog_id"`
	PublishAt time.Time `json:"publish_at"`
}
