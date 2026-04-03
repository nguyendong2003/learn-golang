package repository

import (
	"context"
	"elearning-api/model"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type CartRepository interface {
	AddCourse(ctx context.Context, userID, courseID uuid.UUID) error
	RemoveCourse(ctx context.Context, userID, courseID uuid.UUID) error
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.CartItem, error)
	ClearByUser(ctx context.Context, userID uuid.UUID) error
}

type cartRepository struct {
	db DbRepository
}

func NewCartRepository(db DbRepository) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) AddCourse(ctx context.Context, userID, courseID uuid.UUID) error {
	item := &model.CartItem{UserID: userID, CourseID: courseID}
	return r.db.GetDB().WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(item).Error
}

func (r *cartRepository) RemoveCourse(ctx context.Context, userID, courseID uuid.UUID) error {
	return r.db.GetDB().WithContext(ctx).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Delete(&model.CartItem{}).Error
}

func (r *cartRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.CartItem, error) {
	var items []*model.CartItem
	err := r.db.GetDB().WithContext(ctx).
		Model(&model.CartItem{}).
		Where("user_id = ?", userID).
		Preload("Course").
		Find(&items).Error
	return items, err
}

func (r *cartRepository) ClearByUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.GetDB().WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.CartItem{}).Error
}
