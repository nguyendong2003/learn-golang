package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository[T any] interface {
	GetAll(ctx context.Context) ([]T, error)
	GetWithPagination(ctx context.Context, limit, page int) ([]T, int64, error)
	GetWithPaginationAndFilter(ctx context.Context, limit, page int, filters map[string]any, allowedFields map[string]bool) ([]T, int64, error)
	Filter(ctx context.Context, filters map[string]any, allowedFields map[string]bool) ([]T, error)
	GetByID(ctx context.Context, id uuid.UUID) (*T, error)
	Create(ctx context.Context, data *T) (*T, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]any) (*T, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository[T any] struct {
	db DbRepository
}

func NewBaseRepository[T any](db DbRepository) *repository[T] {
	return &repository[T]{db: db}
}

func (r *repository[T]) GetAll(ctx context.Context) ([]T, error) {
	var data []T

	if err := r.db.GetDB().WithContext(ctx).Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *repository[T]) GetWithPagination(ctx context.Context, limit, page int) ([]T, int64, error) {
	var (
		data  []T
		total int64
	)

	offset := (page - 1) * limit
	db := r.db.GetDB().WithContext(ctx).Model(new(T))

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (r *repository[T]) Filter(
	ctx context.Context,
	filters map[string]any,
	allowedFields map[string]bool,
) ([]T, error) {

	var data []T
	db := r.db.GetDB().WithContext(ctx)

	for key, value := range filters {
		if allowedFields[key] {
			db = db.Where(key+" = ?", value)
		}
	}

	if err := db.Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *repository[T]) GetWithPaginationAndFilter(
	ctx context.Context,
	limit, page int,
	filters map[string]any,
	allowedFields map[string]bool,
) ([]T, int64, error) {

	var (
		data  []T
		total int64
	)

	offset := (page - 1) * limit
	db := r.db.GetDB().WithContext(ctx).Model(new(T))

	for key, value := range filters {
		if allowedFields[key] {
			db = db.Where(key+" = ?", value)
		}
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (r *repository[T]) GetByID(ctx context.Context, id uuid.UUID) (*T, error) {
	var data T

	err := r.db.GetDB().WithContext(ctx).
		First(&data, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *repository[T]) Create(ctx context.Context, data *T) (*T, error) {
	if err := r.db.GetDB().WithContext(ctx).Create(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *repository[T]) Update(ctx context.Context, id uuid.UUID, updates map[string]any) (*T, error) {
	var entity T
	db := r.db.GetDB().WithContext(ctx)

	if err := db.First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&entity).Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := db.First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *repository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.GetDB().WithContext(ctx).
		Delete(new(T), "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
