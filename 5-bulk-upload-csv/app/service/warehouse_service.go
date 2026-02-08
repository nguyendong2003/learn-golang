package service

import (
	"bulk-upload-csv/dto"
	"bulk-upload-csv/interfaces"
	"bulk-upload-csv/model"
	"context"
)

type WarehouseService struct {
	WarehouseRepository interfaces.WarehouseRepositoryInterface
}

func NewWarehouseService(
	WarehouseRepository interfaces.WarehouseRepositoryInterface,
) interfaces.WarehouseServiceInterface {
	return &WarehouseService{
		WarehouseRepository: WarehouseRepository,
	}
}

func (s *WarehouseService) GetList(ctx context.Context, params dto.GetListWarehouseRequest) ([]dto.WarehouseResponse, []error) {
	categories, errs := s.WarehouseRepository.GetList(ctx, params)

	if errs != nil {
		return nil, errs
	}

	response := dto.NewListWarehouseResponse(categories)
	return response, nil
}

func (s *WarehouseService) GetDetail(ctx context.Context, params dto.GetWarehouseDetailRequest) (*dto.WarehouseResponse, []error) {
	Warehouse, errs := s.WarehouseRepository.GetDetail(ctx, model.GetDetailWarehouseParams{
		Id: &params.Id,
	})

	if errs != nil {
		return nil, errs
	}

	response := dto.NewWarehouseDetailResponse(*Warehouse)

	return &response, nil
}

func (s *WarehouseService) Create(ctx context.Context, data dto.CreateWarehouseRequest) (*dto.WarehouseResponse, []error) {
	Warehouse, errs := s.WarehouseRepository.Create(ctx, model.Warehouse{
		Code: data.Code,
	})

	if errs != nil {
		return nil, errs
	}

	response := dto.NewWarehouseDetailResponse(*Warehouse)

	return &response, nil
}
