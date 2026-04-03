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

type CategoryService interface {
	Create(ctx context.Context, request dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	Update(ctx context.Context, id uuid.UUID, request dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*dto.CategoryResponse, error)
	GetAll(ctx context.Context) ([]*dto.CategoryResponse, error)
	GetList(ctx context.Context, request dto.ListCategoryRequest) ([]*dto.CategoryResponse, int64, error)
}

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepository: categoryRepository,
	}
}

func (s *categoryService) Create(ctx context.Context, request dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepository.GetByName(ctx, request.Name)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to check existing category")
	}

	if category != nil {
		validationErrors := map[string]string{}

		if category.Name == request.Name {
			validationErrors["name"] = "Name already exists"
		}

		return nil, apperror.NewValidationError(validationErrors)
	}

	newCategory := &model.Category{
		Name:        request.Name,
		Description: request.Description,
	}

	createdCategory, err := s.categoryRepository.Create(ctx, newCategory)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create category")
	}

	return dto.NewCategoryDetailResponse(createdCategory), nil
}

func (s *categoryService) Update(ctx context.Context, id uuid.UUID, request dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepository.FindByID(ctx, id, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve category")
	}

	if category == nil {
		return nil, apperror.NewNotFoundError("Category not found")
	}

	// Validation
	validationErrors := map[string]string{}

	if category.Name == request.Name && category.ID != id {
		validationErrors["name"] = "Name already exists"
	}

	if len(validationErrors) > 0 {
		return nil, apperror.NewValidationError(validationErrors)
	}

	// Update fields
	category.Name = request.Name
	category.Description = request.Description

	updatedCategory, err := s.categoryRepository.Updates(ctx, category)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to update category")
	}

	return dto.NewCategoryDetailResponse(updatedCategory), nil
}

func (s *categoryService) Delete(ctx context.Context, id uuid.UUID) error {
	category, err := s.categoryRepository.FindByID(ctx, id, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to retrieve category")
	}

	if category == nil {
		return apperror.NewNotFoundError("Category not found")
	}

	if err := s.categoryRepository.Delete(ctx, id); err != nil {
		return apperror.NewInternalServerError("Failed to delete category")
	}

	return nil
}

func (s *categoryService) GetByID(ctx context.Context, id uuid.UUID) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepository.FindByID(ctx, id, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve category")
	}

	if category == nil {
		return nil, apperror.NewNotFoundError("Category not found")
	}

	return dto.NewCategoryDetailResponse(category), nil
}

func (s *categoryService) GetAll(ctx context.Context) ([]*dto.CategoryResponse, error) {
	categories, err := s.categoryRepository.FindAll(ctx, "", nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to retrieve categories")
	}

	return dto.NewListCategoryResponse(categories), nil
}

func (s *categoryService) GetList(ctx context.Context, request dto.ListCategoryRequest) ([]*dto.CategoryResponse, int64, error) {
	// Build sort query
	orderQuery := buildCategorySortQuery(request.SortBy, request.SortOrder)

	// Build filter query
	query, args := buildCategoryQuery(request)

	categories, total, err := s.categoryRepository.List(
		ctx,
		request.Limit,
		request.Offset,
		orderQuery,
		query,
		nil,
		args...,
	)

	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to retrieve categories")
	}

	return dto.NewListCategoryResponse(categories), total, nil
}

func buildCategorySortQuery(sortBy string, sortOrder string) string {
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

func buildCategoryQuery(request dto.ListCategoryRequest) (string, []any) {
	var conditions []string
	var args []any

	filters := map[string]*string{
		"name": request.Name,
	}

	for column, value := range filters {
		util.AddILIKECondition(&conditions, &args, column, value)
	}

	query := strings.Join(conditions, " AND ")

	return query, args
}
