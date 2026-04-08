package service

import (
	"context"
	"strings"
	"time"

	"elearning-api/apperror"
	"elearning-api/consts"
	"elearning-api/dto"
	"elearning-api/job"
	"elearning-api/model"
	"elearning-api/repository"
	"elearning-api/util"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type CourseService interface {
	Create(ctx context.Context, userID uuid.UUID, request dto.CreateCourseRequest) (*dto.CourseResponse, error)
	Update(ctx context.Context, userID, courseID uuid.UUID, request dto.UpdateCourseRequest) (*dto.CourseResponse, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status consts.CourseStatus) (*dto.CourseResponse, error)
	SubmitForReview(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (*dto.CourseResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*dto.CourseResponse, error)
	GetBySlug(ctx context.Context, slug string, userID uuid.UUID, userRole string) (*dto.CourseResponse, error)
	GetList(ctx context.Context, request dto.ListCourseRequest) ([]*dto.CourseResponse, int64, error)

	CreateEvent(ctx context.Context, userID, courseID uuid.UUID, request dto.CourseEventRequest) (*dto.CourseEventResponse, error)
	GetEvents(ctx context.Context, courseID uuid.UUID) ([]*dto.CourseEventResponse, error)
	DeleteEvent(ctx context.Context, userID, courseID, eventID uuid.UUID) error
	UpdateEvent(ctx context.Context, userID, courseID, eventID uuid.UUID, request dto.CourseEventRequest) (*dto.CourseEventResponse, error)

	GetStatistics(ctx context.Context) (*dto.CourseStatisticsResponse, error)
	GetInstructorTaughtCourseRevenue(ctx context.Context, userID uuid.UUID, request dto.PagingRequest) ([]*dto.InstructorTaughtCourseRevenueResponse, int64, error)
	GetNewCoursesLast30Days(ctx context.Context) (*dto.NewCoursesLast30DaysResponse, error)

	GetCoursePurchaseBySessionID(ctx context.Context, sessionID string) (*dto.CoursePurchaseResponse, error)
}

type courseService struct {
	courseRepository         repository.CourseRepository
	coursePurchaseRepository repository.CoursePurchaseRepository
	categoryRepository       repository.CategoryRepository
	instructorProfileService InstructorProfileService
	courseEventRepository    repository.CourseEventRepository
	enrollmentRepository     repository.EnrollmentRepository
	userRepository           repository.UserRepository
	db                       repository.DbRepository
	uploadService            UploadService
	asynqClient              *asynq.Client
	enrollmentService        EnrollmentService
}

func NewCourseService(
	courseRepository repository.CourseRepository,
	coursePurchaseRepository repository.CoursePurchaseRepository,
	categoryRepository repository.CategoryRepository,
	instructorProfileService InstructorProfileService,
	courseEventRepository repository.CourseEventRepository,
	enrollmentRepository repository.EnrollmentRepository,
	userRepository repository.UserRepository,
	db repository.DbRepository,
	uploadService UploadService,
	asynqClient *asynq.Client,
	enrollmentService EnrollmentService,
) CourseService {
	return &courseService{
		courseRepository:         courseRepository,
		coursePurchaseRepository: coursePurchaseRepository,
		categoryRepository:       categoryRepository,
		instructorProfileService: instructorProfileService,
		courseEventRepository:    courseEventRepository,
		enrollmentRepository:     enrollmentRepository,
		userRepository:           userRepository,
		db:                       db,
		uploadService:            uploadService,
		asynqClient:              asynqClient,
		enrollmentService:        enrollmentService,
	}
}

func (s *courseService) Create(ctx context.Context, userID uuid.UUID, request dto.CreateCourseRequest) (*dto.CourseResponse, error) {
	category, err := s.categoryRepository.FindByID(ctx, request.CategoryID, nil)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, apperror.NewNotFoundError("Category not found")
	}
	user, err := s.userRepository.FindByID(ctx, userID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve user")
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("User not found")
	}

	var createdCourse *model.Course
	err = s.db.Transaction(ctx, func(txDb repository.DbRepository) error {
		txCourseRepository := repository.NewCourseRepository(txDb)

		if err := s.validateAndConfirmCourseImageURL(ctx, txDb, request.ImageURL); err != nil {
			return err
		}

		newCourse := &model.Course{
			Title:       request.Title,
			Description: request.Description,
			Image:       request.ImageURL,
			Slug:        util.GenerateSlug(request.Title),
			Price:       request.Price,
			Status:      consts.CourseDraft,
			CategoryID:  request.CategoryID,
			UserID:      userID,
		}

		var err error
		createdCourse, err = txCourseRepository.Create(ctx, newCourse)
		if err != nil {
			return apperror.NewInternalServerError("Failed to create course")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dto.NewCourseDetailResponse(createdCourse), nil
}

func (s *courseService) Update(ctx context.Context, userID, courseID uuid.UUID, request dto.UpdateCourseRequest) (*dto.CourseResponse, error) {
	course, err := s.courseRepository.FindByID(ctx, courseID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve course")
	}
	if course == nil {
		return nil, apperror.NewNotFoundError("Course not found")
	}
	if course.UserID != userID {
		return nil, apperror.NewForbiddenError("You do not have permission to update this course")
	}

	var updatedCourse *model.Course
	err = s.db.Transaction(ctx, func(txDb repository.DbRepository) error {
		txCourseRepository := repository.NewCourseRepository(txDb)

		if request.Title != nil {
			course.Title = *request.Title
		}

		if request.Description != nil {
			course.Description = *request.Description
		}

		if request.ImageURL != nil && *request.ImageURL != "" && course.Image != *request.ImageURL {
			if err := s.validateAndConfirmCourseImageURL(ctx, txDb, *request.ImageURL); err != nil {
				return err
			}
			course.Image = *request.ImageURL
		}

		if request.Price != nil {
			course.Price = *request.Price
		}

		if request.CategoryID != nil {
			categoryID, err := uuid.Parse(*request.CategoryID)
			if err != nil {
				return apperror.NewValidationError(map[string]string{"category_id": "Invalid UUID"})
			}
			course.CategoryID = categoryID
		}

		var err error
		updatedCourse, err = txCourseRepository.Updates(ctx, course)
		if err != nil {
			return apperror.NewInternalServerError("Failed to update course")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dto.NewCourseDetailResponse(updatedCourse), nil
}

func (s *courseService) UpdateStatus(ctx context.Context, id uuid.UUID, status consts.CourseStatus) (*dto.CourseResponse, error) {
	course, err := s.courseRepository.FindByID(ctx, id, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve course")
	}

	if course == nil {
		return nil, apperror.NewNotFoundError("Course not found")
	}

	course.Status = status

	updatedCourse, err := s.courseRepository.Updates(ctx, course)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to update course status")
	}

	return dto.NewCourseDetailResponse(updatedCourse), nil
}

func (s *courseService) SubmitForReview(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (*dto.CourseResponse, error) {
	course, err := s.courseRepository.FindByID(ctx, courseID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve course")
	}

	if course == nil {
		return nil, apperror.NewNotFoundError("Course not found")
	}

	if course.UserID != userID {
		return nil, apperror.NewForbiddenError("You do not have permission to submit this course")
	}

	course.Status = consts.CoursePending

	updatedCourse, err := s.courseRepository.Updates(ctx, course)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to submit course for review")
	}

	return dto.NewCourseDetailResponse(updatedCourse), nil
}

func (s *courseService) Delete(ctx context.Context, id uuid.UUID) error {
	course, err := s.courseRepository.FindByID(ctx, id, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to retrieve course")
	}

	if course == nil {
		return apperror.NewNotFoundError("Course not found")
	}

	if err := s.courseRepository.Delete(ctx, id); err != nil {
		return apperror.NewInternalServerError("Failed to delete course")
	}

	return nil
}

func (s *courseService) GetByID(ctx context.Context, id uuid.UUID) (*dto.CourseResponse, error) {
	course, err := s.courseRepository.FindByID(ctx, id, []repository.Preload{
		repository.PreloadPath(repository.User, repository.InstructorProfile, repository.User),
		repository.PreloadPath(repository.Category),
	})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve course")
	}

	if course == nil {
		return nil, apperror.NewNotFoundError("Course not found")
	}

	return dto.NewCourseDetailResponse(course), nil
}

func (s *courseService) GetBySlug(ctx context.Context, slug string, userID uuid.UUID, userRole string) (*dto.CourseResponse, error) {
	course, err := s.courseRepository.Find(ctx, "slug = ? AND status = ?", []repository.Preload{
		repository.PreloadPath(repository.User, repository.InstructorProfile, repository.User),
		repository.PreloadPath(repository.Category),
	}, slug, consts.CoursePublished)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve course")
	}
	if course == nil {
		return nil, apperror.NewNotFoundError("Course not found")
	}

	result := dto.NewCourseDetailResponse(course)
	if userRole == string(consts.RoleStudent) {
		enrollment, err := s.enrollmentService.FindEnrollment(
			ctx,
			userID,
			course.ID,
		)
		if err != nil {
			return result, nil
		}
		result.IsPurchased = enrollment != nil && enrollment.CanceledAt == nil
	}

	return result, nil
}

func (s *courseService) GetList(ctx context.Context, request dto.ListCourseRequest) ([]*dto.CourseResponse, int64, error) {
	// Build sort query
	orderQuery := buildCourseSortQuery(request.SortBy, request.SortOrder)

	// Build filter query
	query, args := buildCourseQuery(request)

	courses, total, err := s.courseRepository.List(
		ctx,
		request.Limit,
		request.Offset,
		orderQuery,
		query,
		[]repository.Preload{
			repository.PreloadPath(repository.User, repository.InstructorProfile, repository.User),
			repository.PreloadPath(repository.Category),
		},
		args...,
	)
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to retrieve courses")
	}

	return dto.NewListCourseResponse(courses), total, nil
}

func (s *courseService) CreateEvent(ctx context.Context, userID, courseID uuid.UUID, request dto.CourseEventRequest) (*dto.CourseEventResponse, error) {
	course, err := s.courseRepository.FindByID(ctx, courseID, []repository.Preload{repository.CourseEvents})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve course")
	}

	if course == nil {
		return nil, apperror.NewNotFoundError("Course not found")
	}
	if userID != course.UserID {
		return nil, apperror.NewForbiddenError("You do not have permission to create event for this course")
	}
	for _, event := range course.CourseEvents {
		if request.StartTime.UTC().Before(event.EndTime.UTC()) &&
			request.EndTime.UTC().After(event.StartTime.UTC()) {
			return nil, apperror.NewValidationError(map[string]string{"time": "Event time overlaps with existing event"})
		}
	}
	if len(request.AttendeeEmails) > 0 {
		users, err := s.userRepository.FindAll(ctx, "email IN ?", nil, request.AttendeeEmails)
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to validate attendee emails")
		}

		emailToUserID := make(map[string]uuid.UUID, len(users))
		for _, u := range users {
			emailToUserID[u.Email] = u.ID
		}
		var (
			invalid []string
			userIDs []uuid.UUID
		)
		for _, email := range request.AttendeeEmails {
			id, ok := emailToUserID[email]
			if ok {
				userIDs = append(userIDs, id)
			} else {
				invalid = append(invalid, email)
			}
		}
		if len(invalid) > 0 {
			return nil, apperror.NewBadRequestError("Some attendee emails are not registered: " + strings.Join(invalid, ", "))
		}
		enrolledUsers := make(map[uuid.UUID]bool)
		if len(userIDs) > 0 {
			enrollments, err := s.enrollmentRepository.FindAll(ctx, "course_id = ? AND user_id IN ?", nil, courseID, userIDs)
			if err != nil {
				return nil, apperror.NewInternalServerError("Failed to validate attendee emails")
			}
			for _, e := range enrollments {
				enrolledUsers[e.UserID] = true
			}
		}

		for _, email := range request.AttendeeEmails {
			id, ok := emailToUserID[email]
			if !ok || !enrolledUsers[id] {
				invalid = append(invalid, email)
			}
		}
		if len(invalid) > 0 {
			return nil, apperror.NewBadRequestError("Some attendee emails are not enrolled: " + strings.Join(invalid, ", "))
		}
	} else {
		return nil, apperror.NewValidationError(map[string]string{"attendee_emails": "At least one attendee email is required"})
	}

	event := &model.CourseEvent{
		CourseID:                courseID,
		Name:                    request.Name,
		Description:             request.Description,
		Location:                request.Location,
		RoomToken:               util.GenerateRandomString(32),
		StartTime:               request.StartTime.UTC(),
		EndTime:                 request.EndTime.UTC(),
		NotificationBeforeStart: request.NotificationBeforeStart,
		AttendeeEmails:          request.AttendeeEmails,
	}

	event, err = s.courseEventRepository.Create(ctx, event)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create course event")
	}

	if s.asynqClient != nil {
		processAt := event.StartTime.UTC().Add(-time.Duration(event.NotificationBeforeStart) * time.Minute)
		task, opts := job.NewEventNotificationTask(event.ID, event.CourseID, processAt)
		_, _ = s.asynqClient.Enqueue(task, opts...)
	}

	return dto.NewCourseEventResponse(event), nil
}

func (s *courseService) GetEvents(ctx context.Context, courseID uuid.UUID) ([]*dto.CourseEventResponse, error) {
	course, err := s.courseRepository.FindByID(ctx, courseID, []repository.Preload{repository.CourseEvents})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve course")
	}

	if course == nil {
		return nil, apperror.NewNotFoundError("Course not found")
	}
	return dto.NewListCourseEventResponse(course.CourseEvents), nil
}

func (s *courseService) DeleteEvent(ctx context.Context, userID, courseID, eventID uuid.UUID) error {
	course, err := s.courseRepository.FindByID(ctx, courseID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to retrieve course")
	}
	if course == nil {
		return apperror.NewNotFoundError("Course not found")
	}
	if userID != course.UserID {
		return apperror.NewForbiddenError("You do not have permission to delete this event")
	}
	if err := s.courseEventRepository.Delete(ctx, eventID); err != nil {
		return apperror.NewInternalServerError("Failed to delete course event")
	}

	return nil
}

func (s *courseService) UpdateEvent(ctx context.Context, userID, courseID, eventID uuid.UUID, request dto.CourseEventRequest) (*dto.CourseEventResponse, error) {
	course, err := s.courseRepository.FindByID(ctx, courseID, []repository.Preload{repository.CourseEvents})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve course")
	}
	if course == nil {
		return nil, apperror.NewNotFoundError("Course not found")
	}
	if userID != course.UserID {
		return nil, apperror.NewForbiddenError("You do not have permission to update this event")
	}
	var event *model.CourseEvent
	for _, e := range course.CourseEvents {
		if e.ID != eventID && request.StartTime.UTC().Before(e.EndTime.UTC()) &&
			request.EndTime.UTC().After(e.StartTime.UTC()) {
			return nil, apperror.NewValidationError(map[string]string{"time": "Event time overlaps with existing event"})
		}
		if e.ID == eventID {
			event = e
		}
	}
	if event == nil {
		return nil, apperror.NewNotFoundError("Event not found")
	}

	if len(request.AttendeeEmails) > 0 {
		users, err := s.userRepository.FindAll(ctx, "email IN ?", nil, request.AttendeeEmails)
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to validate attendee emails")
		}

		emailToUserID := make(map[string]uuid.UUID, len(users))
		for _, u := range users {
			emailToUserID[u.Email] = u.ID
		}

		var userIDs []uuid.UUID
		for _, email := range request.AttendeeEmails {
			if id, ok := emailToUserID[email]; ok {
				userIDs = append(userIDs, id)
			}
		}

		enrolledUsers := make(map[uuid.UUID]bool)
		if len(userIDs) > 0 {
			enrollments, err := s.enrollmentRepository.FindAll(ctx, "course_id = ? AND user_id IN ?", nil, courseID, userIDs)
			if err != nil {
				return nil, apperror.NewInternalServerError("Failed to validate attendee emails")
			}
			for _, e := range enrollments {
				enrolledUsers[e.UserID] = true
			}
		}

		var invalid []string
		for _, email := range request.AttendeeEmails {
			id, ok := emailToUserID[email]
			if !ok || !enrolledUsers[id] {
				invalid = append(invalid, email)
			}
		}
		if len(invalid) > 0 {
			return nil, apperror.NewBadRequestError("Some attendee emails are not enrolled: " + strings.Join(invalid, ", "))
		}
	} else {
		return nil, apperror.NewValidationError(map[string]string{"attendee_emails": "At least one attendee email is required"})
	}
	event.Name = request.Name
	event.Description = request.Description
	event.Location = request.Location
	event.StartTime = request.StartTime.UTC()
	event.EndTime = request.EndTime.UTC()
	event.NotificationBeforeStart = request.NotificationBeforeStart
	event.AttendeeEmails = request.AttendeeEmails

	updatedEvent, err := s.courseEventRepository.Updates(ctx, event)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to update course event")
	}

	if s.asynqClient != nil {
		processAt := updatedEvent.StartTime.UTC().Add(-time.Duration(updatedEvent.NotificationBeforeStart) * time.Minute)
		task, opts := job.NewEventNotificationTask(updatedEvent.ID, updatedEvent.CourseID, processAt)
		_, _ = s.asynqClient.Enqueue(task, opts...)
	}

	return dto.NewCourseEventResponse(updatedEvent), nil
}

func (s *courseService) GetStatistics(ctx context.Context) (*dto.CourseStatisticsResponse, error) {
	courseStats, err := s.courseRepository.GetStatistics(ctx)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve course statistics")
	}

	return &dto.CourseStatisticsResponse{
		TotalCourses:   courseStats.TotalCourses,
		PendingReviews: courseStats.PendingReviews,
		Drafts:         courseStats.Drafts,
		Published:      courseStats.Published,
		Archived:       courseStats.Archived,
	}, nil
}

func (s *courseService) GetInstructorTaughtCourseRevenue(ctx context.Context, userID uuid.UUID, request dto.PagingRequest) ([]*dto.InstructorTaughtCourseRevenueResponse, int64, error) {
	rows, total, err := s.courseRepository.GetInstructorTaughtCourseRevenue(
		ctx,
		userID,
		request.Limit,
		request.Offset,
		request.SortBy,
		request.SortOrder,
	)
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to retrieve instructor taught courses")
	}

	return dto.NewInstructorTaughtCourseRevenueListResponse(rows), total, nil
}

func (s *courseService) GetNewCoursesLast30Days(ctx context.Context) (*dto.NewCoursesLast30DaysResponse, error) {
	since := time.Now().AddDate(0, 0, -30)
	totalNewCourses, err := s.courseRepository.CountCreatedSince(ctx, since)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to count new courses in the last 30 days")
	}

	return &dto.NewCoursesLast30DaysResponse{
		TotalNewCourses: totalNewCourses,
	}, nil
}

func (s *courseService) GetCoursePurchaseBySessionID(ctx context.Context, sessionID string) (*dto.CoursePurchaseResponse, error) {
	purchase, err := s.coursePurchaseRepository.GetByCheckoutSessionID(ctx, sessionID, []repository.Preload{repository.Details})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get course purchase")
	}
	if purchase == nil {
		return nil, apperror.NewNotFoundError("Course purchase not found")
	}
	return dto.NewCoursePurchaseResponse(purchase), nil
}

func buildCourseSortQuery(sortBy string, sortOrder string) string {
	defaultSort := "created_at DESC"

	if sortBy == "" {
		return defaultSort
	}

	if sortOrder == "" {
		sortOrder = "DESC"
	}

	allowedSort := map[string]bool{
		"created_at":    true,
		"updated_at":    true,
		"average_rate":  true,
		"total_student": true,
		"price":         true,
		"old_price":     true,
	}

	if !allowedSort[sortBy] {
		return defaultSort
	}

	return sortBy + " " + strings.ToUpper(sortOrder)
}

func buildCourseQuery(request dto.ListCourseRequest) (string, []any) {
	var conditions []string
	var args []any

	// LIKE search
	util.AddILIKECondition(&conditions, &args, "title", request.Title)

	// EQUAL filters
	util.AddEqualCondition(&conditions, &args, "category_id", request.CategoryID)
	if request.Status != nil {
		util.AddEqualCondition(&conditions, &args, "status", (*string)(request.Status))
	}

	query := strings.Join(conditions, " AND ")
	return query, args
}

func (s *courseService) validateAndConfirmCourseImageURL(ctx context.Context, txDb repository.DbRepository, imageURL string) error {
	if imageURL == "" {
		return nil
	}
	isValid, err := s.uploadService.ValidateImageURL(ctx, imageURL)
	if !isValid {
		return apperror.NewBadRequestError("Image URL is not valid")
	}
	if err != nil {
		return apperror.NewInternalServerError("Failed to validate image URL")
	}

	trackingRepo := repository.NewPresignedUploadTrackingRepository(txDb)
	presignURL, err := trackingRepo.Find(ctx, "object_url = ?", nil, imageURL)
	if err != nil {
		return apperror.NewInternalServerError("Failed to retrieve presigned URL")
	}
	if presignURL == nil {
		return apperror.NewNotFoundError("Presigned URL not found")
	}
	if presignURL.Status == consts.PresignedUploadStatusConfirmed {
		return apperror.NewBadRequestError("Image URL already used")
	}
	if err := trackingRepo.ConfirmByObjectURL(ctx, imageURL); err != nil {
		return apperror.NewInternalServerError("Failed to confirm image URL")
	}

	return nil
}
