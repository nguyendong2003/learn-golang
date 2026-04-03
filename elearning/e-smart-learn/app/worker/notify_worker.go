package worker

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"elearning-api/job"
	"elearning-api/pkg"
	"elearning-api/repository"

	"github.com/hibiken/asynq"
)

type EmailWorkerHandler struct {
	EventRepository repository.CourseEventRepository
	UserRepository  repository.UserRepository
	EmailProvider   pkg.EmailProvider
	AsynqClient     *asynq.Client
}

func NewEmailWorkerHandler(eventRepo repository.CourseEventRepository, userRepository repository.UserRepository, emailProvider pkg.EmailProvider, asynqClient *asynq.Client) *EmailWorkerHandler {
	return &EmailWorkerHandler{
		EventRepository: eventRepo,
		UserRepository:  userRepository,
		EmailProvider:   emailProvider,
		AsynqClient:     asynqClient,
	}
}

func (h *EmailWorkerHandler) HandleEventNotification(ctx context.Context, t *asynq.Task) error {
	var p job.EventNotificationPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	event, err := h.EventRepository.FindByID(ctx, p.EventID, nil)
	if err != nil {
		return err
	}
	if event == nil {
		return nil
	}
	expectedNotifyAt := event.StartTime.UTC().Add(-time.Duration(event.NotificationBeforeStart) * time.Minute)
	if !p.NotifyAt.IsZero() {
		if math.Abs(p.NotifyAt.UTC().Sub(expectedNotifyAt).Seconds()) > 1 {
			return nil
		}
	}
	attendeeEmails := []string(event.AttendeeEmails)
	users, err := h.UserRepository.FindAll(ctx, "email IN ?", nil, attendeeEmails)
	if err != nil {
		return err
	}
	for _, user := range users {
		payload, _ := json.Marshal(job.EmailEventReminderPayload{
			Email:       user.Email,
			EventName:   event.Name,
			StudentName: user.Name,
			StartTime:   event.StartTime,
			RoomToken:   event.RoomToken,
		})
		_, err = h.AsynqClient.Enqueue(asynq.NewTask(job.TypeSendSingleEmail, payload))
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *EmailWorkerHandler) HandleSendEmail(ctx context.Context, t *asynq.Task) error {
	var p job.EmailEventReminderPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	if err := h.EmailProvider.SendEventReminder(p.Email, p.EventName, p.StudentName, p.StartTime, p.RoomToken); err != nil {
		return err
	}
	return nil
}
