package service

import (
	"context"
	"elearning-api/apperror"
	"elearning-api/consts"
	"elearning-api/dto"
	"elearning-api/repository"
	"slices"
	"time"

	"github.com/google/uuid"
)

type EnrollmentService interface {
	FindEnrollment(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (*dto.EnrollmentResponse, error)
	UpdateLearnedLessons(ctx context.Context, userID uuid.UUID, request dto.StudentCourseProgressRequest) (dto.EnrollmentResponse, error)
	GetMyCourses(ctx context.Context, userID uuid.UUID, request dto.CourseProgresssRequest) ([]*dto.MyCourseItem, int64, error)
	EnrollInCourse(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (*dto.EnrollmentResponse, error)
}

type enrollmentService struct {
	enrollmentRepository               repository.EnrollmentRepository
	subscriptionRepository             repository.SubscriptionRepository
	paymentRepository                  repository.PaymentRepository
	subscriptionRevenueShareRepository repository.SubscriptionRevenueShareRepository
}

func NewEnrollmentService(
	enrollmentRepo repository.EnrollmentRepository,
	subscriptionRepo repository.SubscriptionRepository,
	paymentRepo repository.PaymentRepository,
	subscriptionRevenueShareRepo repository.SubscriptionRevenueShareRepository,
) EnrollmentService {
	return &enrollmentService{
		enrollmentRepository:               enrollmentRepo,
		subscriptionRepository:             subscriptionRepo,
		paymentRepository:                  paymentRepo,
		subscriptionRevenueShareRepository: subscriptionRevenueShareRepo,
	}
}

func (s *enrollmentService) FindEnrollment(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (*dto.EnrollmentResponse, error) {
	enrollment, err := s.enrollmentRepository.Find(
		ctx,
		"user_id=? AND course_id=?",
		nil,
		userID,
		courseID,
	)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to find enrollment")
	}
	if enrollment == nil {
		return nil, apperror.NewNotFoundError("Enrollment not found")
	}
	return dto.NewEnrollmentResponse(enrollment), nil
}

func (s *enrollmentService) UpdateLearnedLessons(ctx context.Context, userID uuid.UUID, request dto.StudentCourseProgressRequest) (dto.EnrollmentResponse, error) {
	userEnrollment, err := s.enrollmentRepository.Find(
		ctx,
		"user_id=? AND course_id=?",
		[]repository.Preload{
			repository.PreloadPath(repository.Course, repository.Chapters, repository.Lessons),
		},
		userID,
		request.CourseId,
	)
	if err != nil {
		return dto.EnrollmentResponse{}, apperror.NewInternalServerError("Failed to find progress")
	}
	if userEnrollment == nil {
		return dto.EnrollmentResponse{}, apperror.NewNotFoundError("User is not enrolled in the course")
	}

	justLearnedLessonId, parseErr := uuid.Parse(request.LessonId)
	if parseErr != nil {
		return dto.EnrollmentResponse{}, apperror.NewBadRequestError("Invalid lesson id")
	}

	isLessonBelongsToCourse := false
	totalLessons := 0
	// Check if the lesson belongs to the course and count total lessons in the course
	for _, chapter := range userEnrollment.Course.Chapters {
		totalLessons += len(chapter.Lessons)
		if isLessonBelongsToCourse {
			continue
		}
		for _, lesson := range chapter.Lessons {
			if lesson.ID == justLearnedLessonId {
				isLessonBelongsToCourse = true
				break
			}
		}
	}

	if !isLessonBelongsToCourse {
		return dto.EnrollmentResponse{}, apperror.NewBadRequestError("Lesson does not belong to the course")
	}

	if slices.Contains([]uuid.UUID(userEnrollment.LearnedLessonIds), justLearnedLessonId) {
		return *dto.NewEnrollmentResponse(userEnrollment), nil
	}
	userEnrollment.LearnedLessonIds = append(userEnrollment.LearnedLessonIds, justLearnedLessonId)
	if _, err = s.enrollmentRepository.AddLearnCourse(
		ctx,
		userEnrollment.UserID,
		userEnrollment.CourseID,
		justLearnedLessonId,
	); err != nil {
		return dto.EnrollmentResponse{}, apperror.NewInternalServerError("Failed to update learned lessons")
	}

	if len(userEnrollment.LearnedLessonIds) == totalLessons {
		if err = s.enrollmentRepository.MarkCourseCompleted(ctx, userEnrollment.UserID, userEnrollment.CourseID); err != nil {
			return dto.EnrollmentResponse{}, apperror.NewInternalServerError("Failed to mark course completed")
		}
	}

	return *dto.NewEnrollmentResponse(userEnrollment), nil
}

func (s *enrollmentService) GetMyCourses(ctx context.Context, userID uuid.UUID, request dto.CourseProgresssRequest) ([]*dto.MyCourseItem, int64, error) {
	enrollments, total, err := s.enrollmentRepository.List(
		ctx,
		request.Limit,
		request.Offset,
		"enrolled_at desc",
		"user_id = ?",
		[]repository.Preload{
			repository.Course,
			repository.PreloadPath(repository.Course, repository.User),
			repository.PreloadPath(repository.Course, repository.Category),
			repository.PreloadPath(repository.Course, repository.Chapters, repository.Lessons),
		},
		userID,
	)
	if err != nil {
		return nil, 0, err
	}

	res := make([]*dto.MyCourseItem, len(enrollments))
	for i, e := range enrollments {
		res[i] = dto.NewMyCourseItem(e)
	}

	return res, total, nil
}

func (s *enrollmentService) EnrollInCourse(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (*dto.EnrollmentResponse, error) {
	enrollment, err := s.enrollmentRepository.EnrollIfNotExists(ctx, userID, courseID, consts.EnrollmentTypeSubscription)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create enrollment")
	}
	if err := s.rebuildCurrentSubscriptionRevenueShares(ctx, userID); err != nil {
		return nil, err
	}
	return dto.NewEnrollmentResponse(enrollment), nil
}

func (s *enrollmentService) rebuildCurrentSubscriptionRevenueShares(ctx context.Context, userID uuid.UUID) error {
	activeSub, err := s.subscriptionRepository.GetActiveByUserID(ctx, userID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to get active subscription")
	}
	if activeSub == nil {
		return nil
	}

	currentPayment, err := s.paymentRepository.GetSucceededPaymentBySubscriptionAndTime(ctx, activeSub.ID, time.Now().UTC(), nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to get current subscription payment")
	}
	if currentPayment == nil {
		return nil
	}

	if err := s.subscriptionRevenueShareRepository.RebuildByPaymentID(ctx, currentPayment.ID); err != nil {
		return apperror.NewInternalServerError("Failed to refresh subscription revenue shares")
	}

	return nil
}
