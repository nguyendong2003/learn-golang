package service

import (
	"bulk-upload-csv/dto"
	"bulk-upload-csv/interfaces"
	"bulk-upload-csv/model"
	"context"
)

type CategoryService struct {
	categoryRepository interfaces.CategoryRepositoryInterface
}

func NewCategoryService(
	categoryRepository interfaces.CategoryRepositoryInterface,
) interfaces.CategoryServiceInterface {
	return &CategoryService{
		categoryRepository: categoryRepository,
	}
}

func (s *CategoryService) GetList(ctx context.Context) ([]dto.CategoryResponse, []error) {
	categories, errs := s.categoryRepository.GetList(ctx)

	if errs != nil {
		return nil, errs
	}

	response := dto.NewListCategoryResponse(categories)
	return response, nil
}

func (s *CategoryService) GetDetail(ctx context.Context, params dto.GetCategoryDetailRequest) (*dto.CategoryResponse, []error) {
	category, errs := s.categoryRepository.GetDetail(ctx, model.GetDetailCategoryParams{
		Id: &params.Id,
	})

	if errs != nil {
		return nil, errs
	}

	response := dto.NewCategoryDetailResponse(*category)

	return &response, nil
}

func (s *CategoryService) Create(ctx context.Context, data dto.CreateCategoryRequest) (*dto.CategoryResponse, []error) {
	category, errs := s.categoryRepository.Create(ctx, model.Category{
		Code: data.Code,
		Name: data.Name,
	})

	if errs != nil {
		return nil, errs
	}

	response := dto.NewCategoryDetailResponse(*category)

	return &response, nil
}
