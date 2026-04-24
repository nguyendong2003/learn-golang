package service

import (
	"context"
	"sort"
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
	Create(ctx context.Context, userID uuid.UUID, userRole string, request dto.CreateCourseRequest) (*dto.CourseResponse, error)
	Update(ctx context.Context, userID, courseID uuid.UUID, request dto.UpdateCourseRequest) (*dto.CourseResponse, error)
	AssignCoupons(ctx context.Context, userID, courseID uuid.UUID, request dto.AssignCourseCouponsRequest) (*dto.CourseResponse, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status consts.CourseStatus) (*dto.CourseResponse, error)
	SubmitForReview(ctx context.Context, userRole string, userID uuid.UUID, courseID uuid.UUID) (*dto.CourseResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*dto.CourseResponse, error)
	GetBySlug(ctx context.Context, slug string, userID uuid.UUID, userRole string) (*dto.CourseResponse, error)
	GetList(ctx context.Context, request dto.ListCourseRequest, userRole string) ([]*dto.CourseResponse, int64, error)

	CreateEvent(ctx context.Context, userID, courseID uuid.UUID, request dto.CourseEventRequest) (*dto.CourseEventResponse, error)
	GetEvents(ctx context.Context, courseID uuid.UUID) ([]*dto.CourseEventResponse, error)
	DeleteEvent(ctx context.Context, userID, courseID, eventID uuid.UUID) error
	UpdateEvent(ctx context.Context, userID, courseID, eventID uuid.UUID, request dto.CourseEventRequest) (*dto.CourseEventResponse, error)

	GetStatistics(ctx context.Context) (*dto.CourseStatisticsResponse, error)
	GetInstructorTaughtCourseRevenue(ctx context.Context, userID uuid.UUID, request dto.PagingRequest) ([]*dto.InstructorTaughtCourseRevenueResponse, int64, error)
	GetNewCoursesLast30Days(ctx context.Context) (*dto.NewCoursesLast30DaysResponse, error)

	GetCoursePurchaseBySessionID(ctx context.Context, sessionID string) (*dto.CoursePurchaseResponse, error)
	GetRecommendedCourses(ctx context.Context, userID uuid.UUID, limit int) ([]*dto.CourseResponse, error)
	GetFeaturedCourses(ctx context.Context, limit int) ([]*dto.CourseResponse, error)
	GetPersonalizedRecommendations(ctx context.Context, userID uuid.UUID, limit int) ([]*dto.CourseResponse, error)
	GetRecommendedCoursesByCategories(ctx context.Context, userID uuid.UUID, limit int) ([]*dto.CourseResponse, error)
}

type courseService struct {
	courseRepository         repository.CourseRepository
	coursePurchaseRepository repository.CoursePurchaseRepository
	couponRepository         repository.CouponRepository
	courseCouponRepository   repository.CourseCouponRepository
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
	couponRepository repository.CouponRepository,
	courseCouponRepository repository.CourseCouponRepository,
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
		couponRepository:         couponRepository,
		courseCouponRepository:   courseCouponRepository,
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

func (s *courseService) Create(ctx context.Context, userID uuid.UUID, userRole string, request dto.CreateCourseRequest) (*dto.CourseResponse, error) {
  if err := validateCreateCourseCouponsRequest(request.Coupons); err != nil {
		return nil, err
	}
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
  // Course status will be CourseDraft by default
  // When creator is Admin, status will be CoursePublished
	createCourseStatus := consts.CourseDraft
	if userRole == string(consts.RoleAdmin) {
		createCourseStatus = consts.CoursePublished
	}
	var createdCourse *model.Course
	err = s.db.Transaction(ctx, func(txDb repository.DbRepository) error {
		txCourseRepository := repository.NewCourseRepository(txDb)
		txCouponRepository := repository.NewCouponRepository(txDb)
		txCourseCouponRepository := repository.NewCourseCouponRepository(txDb)

		if err := s.validateAndConfirmCourseImageURL(ctx, txDb, request.ImageURL); err != nil {
			return err
		}

		newCourse := &model.Course{
			Title:       request.Title,
			Description: request.Description,
			Image:       request.ImageURL,
			Slug:        util.GenerateSlug(request.Title),
			Price:       request.Price,
			Status:      createCourseStatus,
			CategoryID:  request.CategoryID,
			UserID:      userID,
		}

		var err error
		createdCourse, err = txCourseRepository.Create(ctx, newCourse)
		if err != nil {
			return apperror.NewInternalServerError("Failed to create course")
		}

		if len(request.Coupons) > 0 {
			courseCoupons := make([]*model.CourseCoupon, 0, len(request.Coupons))
			for _, reqCoupon := range request.Coupons {
				couponID, parseErr := uuid.Parse(strings.TrimSpace(reqCoupon.CouponID))
				if parseErr != nil {
					return apperror.NewBadRequestError("Invalid coupon id format")
				}

				coupon, err := txCouponRepository.FindByID(ctx, couponID, nil)
				if err != nil {
					return apperror.NewInternalServerError("Failed to retrieve coupon")
				}
				if err := validateCouponAvailability(coupon); err != nil {
					return err
				}
				if coupon.UserID == nil || *coupon.UserID != userID {
					return apperror.NewForbiddenError("Coupon does not belong to this instructor")
				}

				courseCoupons = append(courseCoupons, &model.CourseCoupon{
					CourseID:  createdCourse.ID,
					CouponID:  couponID,
					IsDefault: reqCoupon.IsDefault,
				})
			}

			if err := txCourseCouponRepository.CreateBatch(ctx, courseCoupons); err != nil {
				return apperror.NewInternalServerError("Failed to assign coupons to course")
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dto.NewCourseDetailResponse(createdCourse), nil
}

func validateCreateCourseCouponsRequest(coupons []dto.CreateCourseCouponRequest) error {
	if len(coupons) == 0 {
		return nil
	}

	defaultCount := 0
	seenCouponIDs := make(map[string]bool, len(coupons))

	for _, coupon := range coupons {
		couponID := strings.TrimSpace(coupon.CouponID)
		if couponID == "" {
			return apperror.NewBadRequestError("coupon_id is required")
		}

		if _, exists := seenCouponIDs[couponID]; exists {
			return apperror.NewBadRequestError("Duplicate coupon id in request")
		}
		seenCouponIDs[couponID] = true

		if coupon.IsDefault {
			defaultCount++
			if defaultCount > 1 {
				return apperror.NewBadRequestError("Only one default coupon is allowed")
			}
		}
	}

	return nil
}

func validateUpdateCourseCouponsRequest(addCoupons []dto.CreateCourseCouponRequest, updateCoupons []dto.UpdateCourseCouponRequest, deleteCoupons []string) error {
	if len(addCoupons) == 0 && len(updateCoupons) == 0 && len(deleteCoupons) == 0 {
		return nil
	}

	seenCouponIDs := make(map[string]bool)
	validate := func(couponID string) error {
		couponID = strings.TrimSpace(couponID)
		if couponID == "" {
			return apperror.NewBadRequestError("coupon_id is required")
		}
		if _, exists := seenCouponIDs[couponID]; exists {
			return apperror.NewBadRequestError("Duplicate coupon id in request")
		}
		seenCouponIDs[couponID] = true
		return nil
	}

	for _, addRequest := range addCoupons {
		if err := validate(addRequest.CouponID); err != nil {
			return err
		}
	}

	for _, updateRequest := range updateCoupons {
		if err := validate(updateRequest.CouponID); err != nil {
			return err
		}
	}

	for _, deleteID := range deleteCoupons {
		if err := validate(deleteID); err != nil {
			return err
		}
	}

	return nil
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

	if course.Status == consts.CoursePublished {
		return nil, apperror.NewBadRequestError("Cannot update published course")
	}

	if err := validateUpdateCourseCouponsRequest(request.CouponsAdd, request.CouponsUpdate, request.CouponsDelete); err != nil {
		return nil, err
	}

	var updatedCourse *model.Course
	err = s.db.Transaction(ctx, func(txDb repository.DbRepository) error {
		txCourseRepository := repository.NewCourseRepository(txDb)
		txCouponRepository := repository.NewCouponRepository(txDb)
		txCourseCouponRepository := repository.NewCourseCouponRepository(txDb)

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

		if len(request.CouponsAdd) > 0 || len(request.CouponsUpdate) > 0 || len(request.CouponsDelete) > 0 {
			existingCoupons, err := txCourseCouponRepository.ListByCourseID(ctx, course.ID, nil)
			if err != nil {
				return apperror.NewInternalServerError("Failed to retrieve course coupons")
			}

			// 1. Build current state
			finalMap := make(map[uuid.UUID]*model.CourseCoupon)
			for _, cc := range existingCoupons {
				if cc != nil {
					copy := *cc
					finalMap[cc.CouponID] = &copy
				}
			}

			// 2. DELETE
			for _, idStr := range request.CouponsDelete {
				couponID, err := uuid.Parse(strings.TrimSpace(idStr))
				if err != nil {
					return apperror.NewBadRequestError("Invalid coupon id")
				}
				if _, ok := finalMap[couponID]; !ok {
					return apperror.NewNotFoundError("Course coupon not found")
				}
				delete(finalMap, couponID)
			}

			// 3. UPDATE
			for _, u := range request.CouponsUpdate {
				couponID, err := uuid.Parse(strings.TrimSpace(u.CouponID))
				if err != nil {
					return apperror.NewBadRequestError("Invalid coupon id")
				}

				cc, ok := finalMap[couponID]
				if !ok {
					return apperror.NewNotFoundError("Course coupon not found")
				}

				cc.IsDefault = u.IsDefault
			}

			// 4. ADD
			for _, a := range request.CouponsAdd {
				couponID, err := uuid.Parse(strings.TrimSpace(a.CouponID))
				if err != nil {
					return apperror.NewBadRequestError("Invalid coupon id")
				}

				if _, exists := finalMap[couponID]; exists {
					return apperror.NewBadRequestError("Coupon already assigned")
				}

				coupon, err := txCouponRepository.FindByID(ctx, couponID, nil)
				if err != nil {
					return apperror.NewInternalServerError("Failed to retrieve coupon")
				}
				if err := validateCouponAvailability(coupon); err != nil {
					return err
				}
				if coupon.UserID == nil || *coupon.UserID != userID {
					return apperror.NewForbiddenError("Coupon does not belong to this instructor")
				}

				finalMap[couponID] = &model.CourseCoupon{
					CourseID:  course.ID,
					CouponID:  couponID,
					IsDefault: a.IsDefault,
				}
			}

			// 5. VALIDATE: only 1 default
			defaultCount := 0
			for _, cc := range finalMap {
				if cc.IsDefault {
					defaultCount++
				}
			}
			if defaultCount > 1 {
				return apperror.NewBadRequestError("Only one default coupon is allowed")
			}

			// 6. SYNC DB
			existingMap := make(map[uuid.UUID]*model.CourseCoupon)
			for _, cc := range existingCoupons {
				existingMap[cc.CouponID] = cc
			}

			// DELETE
			for couponID := range existingMap {
				if _, ok := finalMap[couponID]; !ok {
					if err := txCourseCouponRepository.DeleteByCouponID(ctx, course.ID, couponID); err != nil {
						return apperror.NewInternalServerError("Failed to delete coupon")
					}
				}
			}

			// UPSERT
			for couponID, final := range finalMap {
				if existing, ok := existingMap[couponID]; ok {
					if existing.IsDefault != final.IsDefault {
						if _, err := txCourseCouponRepository.Update(ctx, existing.ID, map[string]any{
							"is_default": final.IsDefault,
						}); err != nil {
							return apperror.NewInternalServerError("Failed to update coupon")
						}
					}
				} else {
					if _, err := txCourseCouponRepository.Create(ctx, final); err != nil {
						return apperror.NewInternalServerError("Failed to create coupon")
					}
				}
			}
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

func (s *courseService) AssignCoupons(ctx context.Context, userID, courseID uuid.UUID, request dto.AssignCourseCouponsRequest) (*dto.CourseResponse, error) {
	course, err := s.courseRepository.FindByID(ctx, courseID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve course")
	}
	if course == nil {
		return nil, apperror.NewNotFoundError("Course not found")
	}
	if course.UserID != userID {
		return nil, apperror.NewForbiddenError("You do not have permission to assign coupons to this course")
	}

	if err := validateUpdateCourseCouponsRequest(request.CouponsAdd, request.CouponsUpdate, request.CouponsDelete); err != nil {
		return nil, err
	}

	var updatedCourse *model.Course
	err = s.db.Transaction(ctx, func(txDb repository.DbRepository) error {
		txCouponRepository := repository.NewCouponRepository(txDb)
		txCourseCouponRepository := repository.NewCourseCouponRepository(txDb)

		if len(request.CouponsAdd) > 0 || len(request.CouponsUpdate) > 0 || len(request.CouponsDelete) > 0 {
			existingCoupons, err := txCourseCouponRepository.ListByCourseID(ctx, course.ID, nil)
			if err != nil {
				return apperror.NewInternalServerError("Failed to retrieve course coupons")
			}

			// 1. Build current state
			finalMap := make(map[uuid.UUID]*model.CourseCoupon)
			for _, cc := range existingCoupons {
				if cc != nil {
					copy := *cc
					finalMap[cc.CouponID] = &copy
				}
			}

			// 2. DELETE
			for _, idStr := range request.CouponsDelete {
				couponID, err := uuid.Parse(strings.TrimSpace(idStr))
				if err != nil {
					return apperror.NewBadRequestError("Invalid coupon id")
				}
				if _, ok := finalMap[couponID]; !ok {
					return apperror.NewNotFoundError("Course coupon not found")
				}
				delete(finalMap, couponID)
			}

			// 3. UPDATE
			for _, u := range request.CouponsUpdate {
				couponID, err := uuid.Parse(strings.TrimSpace(u.CouponID))
				if err != nil {
					return apperror.NewBadRequestError("Invalid coupon id")
				}

				cc, ok := finalMap[couponID]
				if !ok {
					return apperror.NewNotFoundError("Course coupon not found")
				}

				cc.IsDefault = u.IsDefault
			}

			// 4. ADD
			for _, a := range request.CouponsAdd {
				couponID, err := uuid.Parse(strings.TrimSpace(a.CouponID))
				if err != nil {
					return apperror.NewBadRequestError("Invalid coupon id")
				}

				if _, exists := finalMap[couponID]; exists {
					return apperror.NewBadRequestError("Coupon already assigned")
				}

				coupon, err := txCouponRepository.FindByID(ctx, couponID, nil)
				if err != nil {
					return apperror.NewInternalServerError("Failed to retrieve coupon")
				}
				if err := validateCouponAvailability(coupon); err != nil {
					return err
				}
				if coupon.UserID == nil || *coupon.UserID != userID {
					return apperror.NewForbiddenError("Coupon does not belong to this instructor")
				}

				finalMap[couponID] = &model.CourseCoupon{
					CourseID:  course.ID,
					CouponID:  couponID,
					IsDefault: a.IsDefault,
				}
			}

			// 5. VALIDATE: only 1 default
			defaultCount := 0
			for _, cc := range finalMap {
				if cc.IsDefault {
					defaultCount++
				}
			}
			if defaultCount > 1 {
				return apperror.NewBadRequestError("Only one default coupon is allowed")
			}

			// 6. SYNC DB
			existingMap := make(map[uuid.UUID]*model.CourseCoupon)
			for _, cc := range existingCoupons {
				existingMap[cc.CouponID] = cc
			}

			// Delete removed coupons
			for couponID := range existingMap {
				if _, ok := finalMap[couponID]; !ok {
					if err := txCourseCouponRepository.DeleteByCouponID(ctx, course.ID, couponID); err != nil {
						return apperror.NewInternalServerError("Failed to delete course coupon")
					}
				}
			}

			// Upsert remaining coupons
			for couponID, cc := range finalMap {
				if existing, ok := existingMap[couponID]; ok && existing != nil {
					// Update existing
					if _, err := txCourseCouponRepository.Updates(ctx, cc); err != nil {
						return apperror.NewInternalServerError("Failed to update course coupon")
					}
				} else {
					// Insert new
					if _, err := txCourseCouponRepository.Create(ctx, cc); err != nil {
						return apperror.NewInternalServerError("Failed to create course coupon")
					}
				}
			}
		}

		updatedCourse = course
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

func (s *courseService) SubmitForReview(ctx context.Context, userRole string, userID uuid.UUID, courseID uuid.UUID) (*dto.CourseResponse, error) {
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
	if userRole == string(consts.RoleAdmin) {
		course.Status = consts.CoursePublished
	}

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

func (s *courseService) GetList(ctx context.Context, request dto.ListCourseRequest, userRole string) ([]*dto.CourseResponse, int64, error) {
	// Build sort query
	orderQuery := buildCourseSortQuery(request.SortBy, request.SortOrder)

	// Build filter query
	query, args := buildCourseQuery(request, userRole)

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

func (s *courseService) GetRecommendedCourses(ctx context.Context, userID uuid.UUID, limit int) ([]*dto.CourseResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	var exclude []uuid.UUID
	if userID != uuid.Nil {
		enrollments, _, err := s.enrollmentRepository.ListByUser(ctx, userID, 0, 0)
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to retrieve enrollments")
		}
		for _, e := range enrollments {
			exclude = append(exclude, e.CourseID)
		}
	}

	courses, err := s.courseRepository.GetTopCoursesByEnrollment(ctx, limit, exclude)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve recommended courses")
	}

	return dto.NewListCourseResponse(courses), nil
}

func (s *courseService) GetFeaturedCourses(ctx context.Context, limit int) ([]*dto.CourseResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	categoryID, err := s.courseRepository.GetTopCategoryID(ctx)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to determine top category")
	}
	if categoryID == uuid.Nil {
		return []*dto.CourseResponse{}, nil
	}

	courses, err := s.courseRepository.GetTopCoursesByCategory(ctx, categoryID, limit, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve featured courses")
	}

	return dto.NewListCourseResponse(courses), nil
}

func (s *courseService) GetPersonalizedRecommendations(ctx context.Context, userID uuid.UUID, limit int) ([]*dto.CourseResponse, error) {
	if userID == uuid.Nil {
		return nil, apperror.NewUnauthorizedError("Unauthorized")
	}
	if limit <= 0 {
		limit = 10
	}

	// Determine the category the user has enrolled in the most
	categoryID, err := s.enrollmentRepository.GetTopCategoryByUser(ctx, userID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to determine user's top category")
	}
	if categoryID == uuid.Nil {
		return []*dto.CourseResponse{}, nil
	}

	// Build exclusion list of user's enrolled courses
	enrollments, _, err := s.enrollmentRepository.ListByUser(ctx, userID, 0, 0)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve enrollments")
	}
	var exclude []uuid.UUID
	for _, e := range enrollments {
		exclude = append(exclude, e.CourseID)
	}

	courses, err := s.courseRepository.GetTopCoursesByCategory(ctx, categoryID, limit, exclude)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve personalized courses")
	}

	return dto.NewListCourseResponse(courses), nil
}

func (s *courseService) GetRecommendedCoursesByCategories(ctx context.Context, userID uuid.UUID, limit int) ([]*dto.CourseResponse, error) {
	if userID == uuid.Nil {
		return nil, apperror.NewUnauthorizedError("Unauthorized")
	}
	if limit <= 0 {
		limit = 10
	}

	// Get categories the user is enrolled in
	categoryIDs, err := s.enrollmentRepository.GetEnrolledCategoryIDs(ctx, userID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve enrolled categories")
	}
	if len(categoryIDs) == 0 {
		return []*dto.CourseResponse{}, nil
	}

	// Build exclusion list of user's enrolled courses
	enrollments, _, err := s.enrollmentRepository.ListByUser(ctx, userID, 0, 0)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve enrollments")
	}
	var exclude []uuid.UUID
	for _, e := range enrollments {
		exclude = append(exclude, e.CourseID)
	}

	type candidate struct {
		CourseID uuid.UUID
		Total    int64
	}

	var candidates []candidate

	for _, cid := range categoryIDs {
		rows, err := s.courseRepository.GetTopCourseIDsAndCountsByCategory(ctx, cid, 1, exclude)
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to retrieve top courses for categories")
		}
		if len(rows) == 0 {
			continue
		}
		candidates = append(candidates, candidate{CourseID: rows[0].CourseID, Total: rows[0].Total})
	}

	if len(candidates) == 0 {
		return []*dto.CourseResponse{}, nil
	}

	sort.Slice(candidates, func(i, j int) bool { return candidates[i].Total > candidates[j].Total })

	if limit > 0 && len(candidates) > limit {
		candidates = candidates[:limit]
	}

	ids := make([]uuid.UUID, 0, len(candidates))
	for _, c := range candidates {
		ids = append(ids, c.CourseID)
	}

	preloads := []repository.Preload{
		repository.PreloadPath(repository.User, repository.InstructorProfile, repository.User),
		repository.PreloadPath(repository.Category),
	}

	courses, err := s.courseRepository.FindAll(ctx, "id IN ?", preloads, ids)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve courses")
	}

	// preserve sorted order
	courseMap := make(map[uuid.UUID]*model.Course, len(courses))
	for _, c := range courses {
		courseMap[c.ID] = c
	}

	var ordered []*model.Course
	for _, cand := range candidates {
		if c, ok := courseMap[cand.CourseID]; ok {
			ordered = append(ordered, c)
		}
	}

	return dto.NewListCourseResponse(ordered), nil
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

func buildCourseQuery(request dto.ListCourseRequest, userRole string) (string, []any) {
	var conditions []string
	var args []any

	// LIKE search
	util.AddILIKECondition(&conditions, &args, "title", request.Title)

	// EQUAL filters
	util.AddEqualCondition(&conditions, &args, "category_id", request.CategoryID)
	if request.Status != nil {
		util.AddEqualCondition(&conditions, &args, "status", (*string)(request.Status))
	}
	if userRole == string(consts.RoleInstructor) && request.UserID != nil {
		util.AddEqualCondition(&conditions, &args, "user_id", request.UserID)
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
