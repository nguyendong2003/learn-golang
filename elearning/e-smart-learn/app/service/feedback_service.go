package service

import (
	"context"

	"elearning-api/apperror"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/repository"

	"github.com/google/uuid"
)

type FeedbackService interface {
	CreateFeedback(ctx context.Context, userId uuid.UUID, request dto.CreateFeedbackRequest) (*dto.FeedbackResponse, error)
	GetFeedbacks(ctx context.Context, courseId uuid.UUID, limit, offset int) ([]*dto.FeedbackResponse, error)
	GetFeaturedFeedbacks(ctx context.Context, limit int) ([]*dto.TestimonialResponse, error)
}

type feedbackService struct {
	feedbackRepository   repository.FeedbackRepository
	userRepository       repository.UserRepository
	enrollmentRepository repository.EnrollmentRepository
}

func NewFeedbackService(feedbackRepository repository.FeedbackRepository, userRepository repository.UserRepository, enrollmentRepository repository.EnrollmentRepository) FeedbackService {
	return &feedbackService{
		feedbackRepository:   feedbackRepository,
		userRepository:       userRepository,
		enrollmentRepository: enrollmentRepository,
	}
}

func (s *feedbackService) CreateFeedback(ctx context.Context, userId uuid.UUID, request dto.CreateFeedbackRequest) (*dto.FeedbackResponse, error) {
	user, err := s.userRepository.FindByID(ctx, userId, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve user")
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("User not found")
	}
	isEnrolled, err := s.enrollmentRepository.CheckEnrollment(ctx, userId, request.CourseID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to check enrollment")
	}
	if !isEnrolled {
		return nil, apperror.NewForbiddenError("User is not enrolled in the course")
	}
	feedback := &model.Feedback{
		UserID:   userId,
		CourseID: request.CourseID,
		Rate:     request.Rating,
		Content:  request.Comment,
	}
	feedback, err = s.feedbackRepository.Create(ctx, feedback)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create feedback")
	}
	return dto.NewFeedbackResponse(feedback), nil
}

func (s *feedbackService) GetFeedbacks(ctx context.Context, courseId uuid.UUID, limit, offset int) ([]*dto.FeedbackResponse, error) {
	feedbacks, _, err := s.feedbackRepository.List(
		ctx,
		limit, offset,
		"created_at desc",
		"course_id = ?",
		[]repository.Preload{repository.User},
		courseId,
	)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve feedbacks")
	}
	return dto.NewFeedbackResponses(feedbacks), nil
}

func (s *feedbackService) GetFeaturedFeedbacks(ctx context.Context, limit int) ([]*dto.TestimonialResponse, error) {
	feedbacks, _, err := s.feedbackRepository.List(
		ctx,
		limit, 0,
		"rate desc",
		"rate >= ?",
		[]repository.Preload{repository.User},
		4,
	)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve feedbacks")
	}
	userIDs := make([]uuid.UUID, 0, len(feedbacks))
	seen := make(map[uuid.UUID]bool)
	for _, f := range feedbacks {
		if !seen[f.UserID] {
			userIDs = append(userIDs, f.UserID)
			seen[f.UserID] = true
		}
	}

	reviewCountMap, err := s.feedbackRepository.CountByUserIDs(ctx, userIDs)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve review count")
	}

	res := make([]*dto.TestimonialResponse, 0, len(feedbacks))
	for _, f := range feedbacks {
		item := &dto.TestimonialResponse{
			ID:          f.ID,
			Comment:     f.Content,
			Rating:      f.Rate,
			ReviewCount: reviewCountMap[f.UserID],
		}
		if f.User != nil {
			item.Name = f.User.Name
			item.Image = f.User.Avatar
		}

		res = append(res, item)
	}

	return res, nil
}
