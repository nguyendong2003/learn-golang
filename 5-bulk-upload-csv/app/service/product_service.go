package service

import (
	"bulk-upload-csv/dto"
	"bulk-upload-csv/interfaces"
	"bulk-upload-csv/model"
	"context"

	"github.com/google/uuid"
)

type ProductService struct {
	productRepository interfaces.ProductRepositoryInterface
}

func NewProductService(
	productRepository interfaces.ProductRepositoryInterface,
) interfaces.ProductServiceInterface {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (s *ProductService) GetList(ctx context.Context, params dto.GetListProductRequest) ([]dto.ProductResponse, []error) {
	categories, errs := s.productRepository.GetList(ctx, params)

	if errs != nil {
		return nil, errs
	}

	response := dto.NewListProductResponse(categories)
	return response, nil
}

func (s *ProductService) GetDetail(ctx context.Context, params dto.GetProductDetailRequest) (*dto.ProductResponse, []error) {
	product, errs := s.productRepository.GetDetail(ctx, model.GetDetailProductParams{
		Id: &params.Id,
	})

	if errs != nil {
		return nil, errs
	}

	response := dto.NewProductDetailResponse(*product)

	return &response, nil
}

func (s *ProductService) Create(ctx context.Context, data dto.CreateProductRequest) (*dto.ProductResponse, []error) {
	product, errs := s.productRepository.Create(ctx, model.Product{
		Sku:        data.Sku,
		CategoryID: uuid.MustParse(data.CategoryID),
	})

	if errs != nil {
		return nil, errs
	}

	response := dto.NewProductDetailResponse(*product)

	return &response, nil
}
