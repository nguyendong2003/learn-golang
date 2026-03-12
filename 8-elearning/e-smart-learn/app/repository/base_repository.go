package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository[T any] interface {
	Create(ctx context.Context, data *T) (*T, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]any) (*T, error)
	Updates(ctx context.Context, entity *T) (*T, error)
	Delete(ctx context.Context, id uuid.UUID) error

	FindByID(ctx context.Context, id uuid.UUID, preloads []Preload) (*T, error)
	Find(ctx context.Context, query string, preloads []Preload, args ...any) (*T, error)
	FindAll(ctx context.Context, query string, preloads []Preload, args ...any) ([]*T, error)

	List(ctx context.Context, limit, offset int, order string, query string, preloads []Preload, args ...any) ([]*T, int64, error)
}

type repository[T any] struct {
	db DbRepository
}

func NewBaseRepository[T any](db DbRepository) *repository[T] {
	return &repository[T]{db: db}
}

func (r *repository[T]) baseQuery(ctx context.Context) *gorm.DB {
	return r.db.GetDB().
		WithContext(ctx).
		Model(new(T))
}

func (r *repository[T]) Create(ctx context.Context, data *T) (*T, error) {
	if err := r.baseQuery(ctx).Create(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *repository[T]) Update(
	ctx context.Context,
	id uuid.UUID,
	updates map[string]any,
) (*T, error) {

	entity, err := r.FindByID(ctx, id, nil)
	if err != nil {
		return nil, err
	}

	if entity == nil {
		return nil, gorm.ErrRecordNotFound
	}

	if err := r.baseQuery(ctx).
		Model(entity).
		Updates(updates).Error; err != nil {
		return nil, err
	}

	// reload entity after update
	if err := r.baseQuery(ctx).
		First(entity, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *repository[T]) Updates(
	ctx context.Context,
	entity *T,
) (*T, error) {

	result := r.baseQuery(ctx).
		Model(entity).
		Updates(entity)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return entity, nil
}

func (r *repository[T]) Delete(ctx context.Context, id uuid.UUID) error {

	result := r.baseQuery(ctx).
		Delete(new(T), "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *repository[T]) FindByID(ctx context.Context, id uuid.UUID, preloads []Preload) (*T, error) {
	db := r.baseQuery(ctx)

	if len(preloads) > 0 {
		db = applyPreloads(db, preloads)
	}

	var entity T

	if err := db.Where("id = ?", id).First(&entity).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *repository[T]) Find(
	ctx context.Context,
	query string,
	preloads []Preload,
	args ...any,
) (*T, error) {

	var entity T

	db := r.baseQuery(ctx)
	if query != "" {
		db = db.Where(query, args...)
	}

	if len(preloads) > 0 {
		db = applyPreloads(db, preloads)
	}

	err := db.First(&entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *repository[T]) FindAll(
	ctx context.Context,
	query string,
	preloads []Preload,
	args ...any,
) ([]*T, error) {
	var data []*T

	db := r.baseQuery(ctx)
	if query != "" {
		db = db.Where(query, args...)
	}

	if len(preloads) > 0 {
		db = applyPreloads(db, preloads)
	}

	if err := db.Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (r *repository[T]) List(
	ctx context.Context,
	limit, offset int,
	order string,
	query string,
	preloads []Preload,
	args ...any,
) ([]*T, int64, error) {

	var (
		data  []*T
		total int64
	)

	db := r.baseQuery(ctx)
	if query != "" {
		db = db.Where(query, args...)
	}

	// count total records
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// apply preloads only for fetching data so it won't affect count performance
	if len(preloads) > 0 {
		db = applyPreloads(db, preloads)
	}

	if order != "" {
		db = db.Order(order)
	}

	// apply pagination
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}

	// fetch data
	if err := db.Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}
