package repository

import (
	"context"

	"elearning-api/consts"
	"elearning-api/model"

	"github.com/google/uuid"
)

type CoursePurchaseRepository interface {
	Repository[model.CoursePurchase]
	GetByCheckoutSessionID(ctx context.Context, checkoutSessionID string, preloads []Preload) (*model.CoursePurchase, error)
	ExistsPaidByUserAndCourse(ctx context.Context, userID, courseID uuid.UUID) (bool, error)
	ListPaidByUserID(ctx context.Context, userID uuid.UUID, preloads []Preload) ([]*model.CoursePurchase, error)
}

type coursePurchaseRepository struct {
	*repository[model.CoursePurchase]
}

func NewCoursePurchaseRepository(db DbRepository) CoursePurchaseRepository {
	return &coursePurchaseRepository{
		repository: NewBaseRepository[model.CoursePurchase](db),
	}
}

func (r *coursePurchaseRepository) GetByCheckoutSessionID(ctx context.Context, checkoutSessionID string, preloads []Preload) (*model.CoursePurchase, error) {
	return r.Find(ctx, "stripe_checkout_session_id = ?", preloads, checkoutSessionID)
}

func (r *coursePurchaseRepository) ExistsPaidByUserAndCourse(ctx context.Context, userID, courseID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.GetDB().WithContext(ctx).
		Table("course_purchases cp").
		Joins("JOIN course_purchase_details cpd ON cpd.course_purchase_id = cp.id").
		Where("cp.user_id = ? AND cpd.course_id = ? AND cp.status = ? AND cp.deleted_at IS NULL AND cpd.deleted_at IS NULL", userID, courseID, consts.CoursePurchaseStatusPaid).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *coursePurchaseRepository) ListPaidByUserID(ctx context.Context, userID uuid.UUID, preloads []Preload) ([]*model.CoursePurchase, error) {
	db := r.baseQuery(ctx).
		Where("user_id = ? AND status = ?", userID, consts.CoursePurchaseStatusPaid).
		Order("created_at DESC")

	if len(preloads) > 0 {
		db = applyPreloads(db, preloads)
	}

	var purchases []*model.CoursePurchase
	if err := db.Find(&purchases).Error; err != nil {
		return nil, err
	}

	return purchases, nil
}

type CoursePurchaseDetailRepository interface {
	Repository[model.CoursePurchaseDetail]
	ListByPurchaseID(ctx context.Context, purchaseID uuid.UUID, preloads []Preload) ([]*model.CoursePurchaseDetail, error)
	CreateBatch(ctx context.Context, details []*model.CoursePurchaseDetail) error
}

type coursePurchaseDetailRepository struct {
	*repository[model.CoursePurchaseDetail]
}

func NewCoursePurchaseDetailRepository(db DbRepository) CoursePurchaseDetailRepository {
	return &coursePurchaseDetailRepository{
		repository: NewBaseRepository[model.CoursePurchaseDetail](db),
	}
}

func (r *coursePurchaseDetailRepository) ListByPurchaseID(ctx context.Context, purchaseID uuid.UUID, preloads []Preload) ([]*model.CoursePurchaseDetail, error) {
	return r.FindAll(ctx, "course_purchase_id = ?", preloads, purchaseID)
}

func (r *coursePurchaseDetailRepository) CreateBatch(ctx context.Context, details []*model.CoursePurchaseDetail) error {
	if len(details) == 0 {
		return nil
	}
	return r.db.GetDB().WithContext(ctx).Create(&details).Error
}
