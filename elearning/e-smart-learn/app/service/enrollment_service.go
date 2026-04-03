package service

import (
	"context"
	"elearning-api/apperror"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/repository"
	"slices"

	"github.com/google/uuid"
)

type EnrollmentService interface {
	FindEnrollment(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (*dto.EnrollmentResponse, error)
	UpdateLearnedLessons(ctx context.Context, userID uuid.UUID, request dto.StudentCourseProgressRequest) (dto.EnrollmentResponse, error)
	GetMyCourses(ctx context.Context, userID uuid.UUID, request dto.CourseProgresssRequest) ([]*dto.MyCourseItem, int64, error)
	EnrollInCourse(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (*dto.EnrollmentResponse, error)
}

type enrollmentService struct {
	enrollmentRepository repository.EnrollmentRepository
	paymentRepository    repository.PaymentRepository
}

func NewEnrollmentService(enrollmentRepo repository.EnrollmentRepository, paymentRepo repository.PaymentRepository) EnrollmentService {
	return &enrollmentService{
		enrollmentRepository: enrollmentRepo,
		paymentRepository:    paymentRepo,
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
	enrollment, err := s.enrollmentRepository.Find(ctx, "user_id=? AND course_id=?", nil, userID, courseID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to find enrollment")
	}
	if enrollment != nil {
		return nil, apperror.NewBadRequestError("User is already enrolled in the course")
	}

	newEnrollment := &model.Enrollment{
		UserID:   userID,
		CourseID: courseID,
	}

	if enrollment, err = s.enrollmentRepository.Create(ctx, newEnrollment); err != nil {
		return nil, apperror.NewInternalServerError("Failed to create enrollment")
	}
	return dto.NewEnrollmentResponse(enrollment), nil
}
