package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"elearning-api/consts"
	"elearning-api/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Repository[model.Subscription]
	GetLatestByUserID(ctx context.Context, userID uuid.UUID, preloads []Preload) (*model.Subscription, error)
	GetByCheckoutSessionID(ctx context.Context, checkoutSessionID string, preloads []Preload) (*model.Subscription, error)
	GetByStripeSubscriptionID(ctx context.Context, stripeSubscriptionID string, preloads []Preload) (*model.Subscription, error)
	GetActiveByUserID(ctx context.Context, userID uuid.UUID, preloads []Preload) (*model.Subscription, error)
	ListStripeSyncCandidates(ctx context.Context, preloads []Preload) ([]*model.Subscription, error)
	ListByUserID(ctx context.Context, userID uuid.UUID, preloads []Preload) ([]*model.Subscription, error)
	ListPendingByCheckoutSessionID(ctx context.Context, preloads []Preload) ([]*model.Subscription, error)
	CountByPlanID(ctx context.Context, planID uuid.UUID) (int64, error)
	CountActiveAsOf(ctx context.Context, asOf time.Time) (int64, error)
	CountRetainedAsOf(ctx context.Context, previousAsOf, currentAsOf time.Time) (int64, error)
	ListAutoRenewByPlanID(ctx context.Context, planID uuid.UUID, preloads []Preload) ([]*model.Subscription, error)
}

type subscriptionRepository struct {
	*repository[model.Subscription]
}

func NewSubscriptionRepository(db DbRepository) SubscriptionRepository {
	return &subscriptionRepository{
		repository: NewBaseRepository[model.Subscription](db),
	}
}

func (r *subscriptionRepository) GetLatestByUserID(ctx context.Context, userID uuid.UUID, preloads []Preload) (*model.Subscription, error) {
	db := r.baseQuery(ctx).Where("user_id = ?", userID).Order("created_at DESC")
	if len(preloads) > 0 {
		db = applyPreloads(db, preloads)
	}

	var sub model.Subscription
	if err := db.First(&sub).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &sub, nil
}

func (r *subscriptionRepository) GetByCheckoutSessionID(ctx context.Context, checkoutSessionID string, preloads []Preload) (*model.Subscription, error) {
	return r.Find(ctx, "stripe_checkout_session_id = ?", preloads, checkoutSessionID)
}

func (r *subscriptionRepository) GetByStripeSubscriptionID(ctx context.Context, stripeSubscriptionID string, preloads []Preload) (*model.Subscription, error) {
	return r.Find(ctx, "stripe_subscription_id = ?", preloads, stripeSubscriptionID)
}

func (r *subscriptionRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID, preloads []Preload) (*model.Subscription, error) {
	now := time.Now().UTC()
	query := fmt.Sprintf(`
		user_id = ?
		AND status IN ('%s', '%s', '%s')
		AND (ended_at IS NULL OR ended_at > ?)
	`, consts.SubscriptionStatusActive, consts.SubscriptionStatusTrialing, consts.SubscriptionStatusPastDue)
	db := r.baseQuery(ctx).Where(query, userID, now).Order("started_at DESC, created_at DESC")
	if len(preloads) > 0 {
		db = applyPreloads(db, preloads)
	}

	var sub model.Subscription
	if err := db.First(&sub).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &sub, nil
}

func (r *subscriptionRepository) ListStripeSyncCandidates(ctx context.Context, preloads []Preload) ([]*model.Subscription, error) {
	now := time.Now().UTC()
	db := r.baseQuery(ctx).
		Where("stripe_subscription_id IS NOT NULL AND stripe_subscription_id <> ''").
		Where("status IN (?, ?, ?)",
			consts.SubscriptionStatusActive,
			consts.SubscriptionStatusTrialing,
			consts.SubscriptionStatusPastDue,
		).
		Where("ended_at IS NULL OR ended_at > ?", now).
		Order("started_at ASC, created_at ASC")

	if len(preloads) > 0 {
		db = applyPreloads(db, preloads)
	}

	var subs []*model.Subscription
	if err := db.Find(&subs).Error; err != nil {
		return nil, err
	}

	return subs, nil
}

func (r *subscriptionRepository) ListByUserID(ctx context.Context, userID uuid.UUID, preloads []Preload) ([]*model.Subscription, error) {
	db := r.baseQuery(ctx).Where("user_id = ?", userID).Order("created_at DESC")

	if len(preloads) > 0 {
		db = applyPreloads(db, preloads)
	}

	var subs []*model.Subscription
	if err := db.Find(&subs).Error; err != nil {
		return nil, err
	}

	return subs, nil
}

func (r *subscriptionRepository) ListPendingByCheckoutSessionID(ctx context.Context, preloads []Preload) ([]*model.Subscription, error) {
	db := r.baseQuery(ctx).
		Where("stripe_checkout_session_id IS NOT NULL AND stripe_checkout_session_id <> ''").
		Where("status IN (?, ?, ?, ?)",
			consts.SubscriptionStatusIncomplete,
			consts.SubscriptionStatusPastDue,
			consts.SubscriptionStatusIncompleteExpired,
			consts.SubscriptionStatusUnpaid,
		).
		Order("created_at ASC")

	if len(preloads) > 0 {
		db = applyPreloads(db, preloads)
	}

	var subs []*model.Subscription
	if err := db.Find(&subs).Error; err != nil {
		return nil, err
	}

	return subs, nil
}

func (r *subscriptionRepository) CountByPlanID(ctx context.Context, planID uuid.UUID) (int64, error) {
	var count int64
	err := r.baseQuery(ctx).
		Where("plan_id = ? AND deleted_at IS NULL", planID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *subscriptionRepository) CountActiveAsOf(ctx context.Context, asOf time.Time) (int64, error) {
	query := `
		SELECT COUNT(1)
		FROM subscriptions
		WHERE deleted_at IS NULL
		  AND status IN (?, ?, ?)
		  AND started_at <= ?
		  AND (ended_at IS NULL OR ended_at > ?)
	`

	var count int64
	if err := r.db.GetDB().WithContext(ctx).Raw(
		query,
		consts.SubscriptionStatusActive,
		consts.SubscriptionStatusTrialing,
		consts.SubscriptionStatusPastDue,
		asOf,
		asOf,
	).Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *subscriptionRepository) CountRetainedAsOf(ctx context.Context, previousAsOf, currentAsOf time.Time) (int64, error) {
	query := `
		SELECT COUNT(DISTINCT prev.user_id)
		FROM subscriptions prev
		WHERE prev.deleted_at IS NULL
		  AND prev.status IN (?, ?, ?)
		  AND prev.started_at <= ?
		  AND (prev.ended_at IS NULL OR prev.ended_at > ?)
		  AND EXISTS (
			SELECT 1
			FROM subscriptions curr
			WHERE curr.deleted_at IS NULL
			  AND curr.user_id = prev.user_id
			  AND curr.status IN (?, ?, ?)
			  AND curr.started_at <= ?
			  AND (curr.ended_at IS NULL OR curr.ended_at > ?)
		  )
	`

	var count int64
	if err := r.db.GetDB().WithContext(ctx).Raw(
		query,
		consts.SubscriptionStatusActive,
		consts.SubscriptionStatusTrialing,
		consts.SubscriptionStatusPastDue,
		previousAsOf,
		previousAsOf,
		consts.SubscriptionStatusActive,
		consts.SubscriptionStatusTrialing,
		consts.SubscriptionStatusPastDue,
		currentAsOf,
		currentAsOf,
	).Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *subscriptionRepository) ListAutoRenewByPlanID(ctx context.Context, planID uuid.UUID, preloads []Preload) ([]*model.Subscription, error) {
	query := fmt.Sprintf(`
		plan_id = ?
		AND cancel_at_period_end = false
		AND status IN ('%s', '%s', '%s')
		AND deleted_at IS NULL
	`, consts.SubscriptionStatusActive, consts.SubscriptionStatusTrialing, consts.SubscriptionStatusPastDue)

	return r.FindAll(ctx, query, preloads, planID)
}
