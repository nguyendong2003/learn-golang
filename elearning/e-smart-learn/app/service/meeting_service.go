package service

import (
	"context"

	"elearning-api/apperror"
	"elearning-api/repository"

	"github.com/google/uuid"
)

type MeetingParticipant struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Avatar string    `json:"avatar,omitempty"`
	Role   string    `json:"role"`
}

type MeetingJoinContext struct {
	RoomID      string
	EventID     uuid.UUID
	CourseID    uuid.UUID
	Participant MeetingParticipant
}

type MeetingService interface {
	ValidateJoin(ctx context.Context, userID, eventID uuid.UUID, roomToken string) (*MeetingJoinContext, error)
}

type meetingService struct {
	courseEventRepository repository.CourseEventRepository
	courseRepository      repository.CourseRepository
	userRepository        repository.UserRepository
}

func NewMeetingService(
	courseEventRepository repository.CourseEventRepository,
	courseRepository repository.CourseRepository,
	userRepository repository.UserRepository,
) MeetingService {
	return &meetingService{
		courseEventRepository: courseEventRepository,
		courseRepository:      courseRepository,
		userRepository:        userRepository,
	}
}

func (s *meetingService) ValidateJoin(ctx context.Context, userID, eventID uuid.UUID, roomToken string) (*MeetingJoinContext, error) {
	event, err := s.courseEventRepository.FindByID(ctx, eventID, []repository.Preload{repository.Course})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve event")
	}
	if event == nil {
		return nil, apperror.NewNotFoundError("Event not found")
	}

	if event.RoomToken != roomToken {
		return nil, apperror.NewForbiddenError("Invalid room token")
	}
	user, err := s.userRepository.FindByID(ctx, userID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve user")
	}
	if user == nil {
		return nil, apperror.NewUnauthorizedError("User not found")
	}

	course := event.Course
	role := ""
	if course.UserID == userID {
		role = "teacher"
	} else {
		for _, email := range event.AttendeeEmails {
			if email == user.Email {
				role = "student"
				break
			}
		}
		if role == "" {
			return nil, apperror.NewForbiddenError("You are not allowed to join this meeting")
		}
	}

	displayName := user.Name
	if displayName == "" {
		displayName = user.Username
	}

	return &MeetingJoinContext{
		RoomID:   event.ID.String(),
		EventID:  event.ID,
		CourseID: event.CourseID,
		Participant: MeetingParticipant{
			UserID: user.ID,
			Name:   displayName,
			Avatar: user.Avatar,
			Role:   role,
		},
	}, nil
}
