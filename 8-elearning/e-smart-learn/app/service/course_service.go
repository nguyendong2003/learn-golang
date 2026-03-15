package service

import (
	"context"
	"elearning-api/apperror"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/repository"
	"elearning-api/util"
	"strings"

	"github.com/google/uuid"
)

type CourseService interface {
	Create(ctx context.Context, userID uuid.UUID, request dto.CreateCourseRequest) (*dto.CourseResponse, error)
	Update(ctx context.Context, id uuid.UUID, request dto.UpdateCourseRequest) (*dto.CourseResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*dto.CourseResponse, error)
	GetBySlug(ctx context.Context, slug string) (*dto.CourseResponse, error)
	GetList(ctx context.Context, request dto.ListCourseRequest) ([]*dto.CourseResponse, int64, error)
}

type courseService struct {
	courseRepository         repository.CourseRepository
	categoryService          CategoryService
	instructorProfileService InstructorProfileService
}

func NewCourseService(
	courseRepository repository.CourseRepository,
	categoryService CategoryService,
	instructorProfileService InstructorProfileService) CourseService {
	return &courseService{
		courseRepository:         courseRepository,
		categoryService:          categoryService,
		instructorProfileService: instructorProfileService,
	}
}

func (s *courseService) Create(ctx context.Context, userID uuid.UUID, request dto.CreateCourseRequest) (*dto.CourseResponse, error) {
	// Validate category ID
	categoryID, err := uuid.Parse(request.CategoryID)
	if err != nil {
		return nil, apperror.NewValidationError(map[string]string{"category_id": "Invalid UUID"})
	}

	category, err := s.categoryService.GetByID(ctx, categoryID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve category")
	}

	if category == nil {
		return nil, apperror.NewNotFoundError("Category not found")
	}

	newCourse := &model.Course{
		Title:       request.Title,
		Description: request.Description,
		Slug:        util.GenerateSlug(request.Title),
		Price:       request.Price,
		OldPrice:    request.Price,
		CategoryID:  categoryID,
		UserID:      userID,
	}

	createdCourse, err := s.courseRepository.Create(ctx, newCourse)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create course")
	}

	course, err := s.courseRepository.FindByID(ctx, createdCourse.ID, []repository.Preload{
		repository.PreloadPath(repository.User, repository.InstructorProfile),
		repository.PreloadPath(repository.Category),
	})

	return dto.NewCourseDetailResponse(course), nil
}

func (s *courseService) Update(ctx context.Context, id uuid.UUID, request dto.UpdateCourseRequest) (*dto.CourseResponse, error) {
	course, err := s.courseRepository.FindByID(ctx, id, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve course")
	}

	if course == nil {
		return nil, apperror.NewNotFoundError("Course not found")
	}

	if request.Title != nil {
		course.Title = *request.Title
	}

	if request.Description != nil {
		course.Description = *request.Description
	}

	if request.Price != nil {
		course.Price = *request.Price
		course.OldPrice = *request.Price
	}

	if request.CategoryID != nil {
		categoryID, err := uuid.Parse(*request.CategoryID)
		if err != nil {
			return nil, apperror.NewValidationError(map[string]string{"category_id": "Invalid UUID"})
		}
		course.CategoryID = categoryID
	}

	updatedCourse, err := s.courseRepository.Updates(ctx, course)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to update course")
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
		repository.PreloadPath(repository.InstructorProfile),
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

func (s *courseService) GetBySlug(ctx context.Context, slug string) (*dto.CourseResponse, error) {
	course, err := s.courseRepository.Find(ctx, "slug = ?", []repository.Preload{
		repository.PreloadPath(repository.InstructorProfile),
		repository.PreloadPath(repository.Category),
	}, slug)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve course")
	}

	if course == nil {
		return nil, apperror.NewNotFoundError("Course not found")
	}

	return dto.NewCourseDetailResponse(course), nil
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
			repository.PreloadPath(repository.InstructorProfile),
			repository.PreloadPath(repository.Category),
		},
		args...,
	)

	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to retrieve categories")
	}

	return dto.NewListCourseResponse(courses), total, nil
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
		"created_at": true,
		"updated_at": true,
	}

	if !allowedSort[sortBy] {
		return defaultSort
	}

	return sortBy + " " + strings.ToUpper(sortOrder)
}

func buildCourseQuery(request dto.ListCourseRequest) (string, []any) {
	var conditions []string
	var args []any

	filters := map[string]*string{
		"title":       request.Title,
		"category_id": request.CategoryID,
		"status":      (*string)(request.Status),
	}

	for column, value := range filters {
		util.AddILIKECondition(&conditions, &args, column, value)
	}

	query := strings.Join(conditions, " AND ")

	return query, args
}
