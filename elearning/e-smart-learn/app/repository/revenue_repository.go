package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type RevenueRepository interface {
	GetInstructorRevenue(ctx context.Context, userID uuid.UUID) (*RevenueBreakdownCents, *RevenueBreakdownCents, error)
	GetInstructorRevenueByDay(ctx context.Context, userID uuid.UUID, date time.Time) (*RevenueBreakdownCents, *RevenueBreakdownCents, error)
	GetInstructorRevenueByMonth(ctx context.Context, userID uuid.UUID, year, month int) (*RevenueBreakdownCents, *RevenueBreakdownCents, error)
	GetInstructorRevenueByYear(ctx context.Context, userID uuid.UUID, year int) (*RevenueBreakdownCents, *RevenueBreakdownCents, error)

	GetAdminRevenue(ctx context.Context) (*RevenueBreakdownCents, *RevenueBreakdownCents, error)
	GetAdminRevenueByDay(ctx context.Context, date time.Time) (*RevenueBreakdownCents, *RevenueBreakdownCents, error)
	GetAdminRevenueByMonth(ctx context.Context, year, month int) (*RevenueBreakdownCents, *RevenueBreakdownCents, error)
	GetAdminRevenueByYear(ctx context.Context, year int) (*RevenueBreakdownCents, *RevenueBreakdownCents, error)

	GetAdminTransactions(ctx context.Context, limit, offset int, sortBy, sortOrder string) ([]*RevenueTransactionRow, int64, error)
}

type RevenueBreakdownCents struct {
	TotalAmount     int64 `gorm:"column:total_amount"`
	InstructorGross int64 `gorm:"column:instructor_gross"`
	PlatformGross   int64 `gorm:"column:platform_gross"`
	StripeFee       int64 `gorm:"column:stripe_fee"`
	InstructorNet   int64 `gorm:"column:instructor_net"`
	PlatformNet     int64 `gorm:"column:platform_net"`
}

type RevenueTransactionRow struct {
	TransactionID string    `gorm:"column:transaction_id"`
	UserID        string    `gorm:"column:user_id"`
	UserEmail     string    `gorm:"column:user_email"`
	UserUsername  string    `gorm:"column:user_username"`
	UserName      string    `gorm:"column:user_name"`
	UserAvatar    string    `gorm:"column:user_avatar"`
	Method        string    `gorm:"column:method"`
	Amount        float64   `gorm:"column:amount"`
	Type          string    `gorm:"column:type"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

type revenueRepository struct {
	db DbRepository
}

func NewRevenueRepository(db DbRepository) RevenueRepository {
	return &revenueRepository{db: db}
}

func (r *revenueRepository) GetInstructorRevenue(ctx context.Context, userID uuid.UUID) (*RevenueBreakdownCents, *RevenueBreakdownCents, error) {
	coursePurchase, err := r.getBreakdown(ctx, "SELECT total_amount, instructor_gross, platform_gross, stripe_fee, instructor_net, platform_net FROM get_instructor_course_purchase_revenue(?)", userID)
	if err != nil {
		return nil, nil, err
	}

	subscription, err := r.getBreakdown(ctx, "SELECT total_amount, instructor_gross, platform_gross, stripe_fee, instructor_net, platform_net FROM get_instructor_subscription_revenue(?)", userID)
	if err != nil {
		return nil, nil, err
	}

	return coursePurchase, subscription, nil
}

func (r *revenueRepository) GetInstructorRevenueByDay(ctx context.Context, userID uuid.UUID, date time.Time) (*RevenueBreakdownCents, *RevenueBreakdownCents, error) {
	coursePurchase, err := r.getBreakdown(ctx, "SELECT total_amount, instructor_gross, platform_gross, stripe_fee, instructor_net, platform_net FROM get_instructor_course_purchase_revenue_by_day(?, ?)", userID, date)
	if err != nil {
		return nil, nil, err
	}

	subscription, err := r.getBreakdown(ctx, "SELECT total_amount, instructor_gross, platform_gross, stripe_fee, instructor_net, platform_net FROM get_instructor_subscription_revenue_by_day(?, ?)", userID, date)
	if err != nil {
		return nil, nil, err
	}

	return coursePurchase, subscription, nil
}

func (r *revenueRepository) GetInstructorRevenueByMonth(ctx context.Context, userID uuid.UUID, year, month int) (*RevenueBreakdownCents, *RevenueBreakdownCents, error) {
	coursePurchase, err := r.getBreakdown(ctx, "SELECT total_amount, instructor_gross, platform_gross, stripe_fee, instructor_net, platform_net FROM get_instructor_course_purchase_revenue_by_month(?, ?, ?)", userID, year, month)
	if err != nil {
		return nil, nil, err
	}

	subscription, err := r.getBreakdown(ctx, "SELECT total_amount, instructor_gross, platform_gross, stripe_fee, instructor_net, platform_net FROM get_instructor_subscription_revenue_by_month(?, ?, ?)", userID, year, month)
	if err != nil {
		return nil, nil, err
	}

	return coursePurchase, subscription, nil
}

func (r *revenueRepository) GetInstructorRevenueByYear(ctx context.Context, userID uuid.UUID, year int) (*RevenueBreakdownCents, *RevenueBreakdownCents, error) {
	coursePurchase, err := r.getBreakdown(ctx, "SELECT total_amount, instructor_gross, platform_gross, stripe_fee, instructor_net, platform_net FROM get_instructor_course_purchase_revenue_by_year(?, ?)", userID, year)
	if err != nil {
		return nil, nil, err
	}

	subscription, err := r.getBreakdown(ctx, "SELECT total_amount, instructor_gross, platform_gross, stripe_fee, instructor_net, platform_net FROM get_instructor_subscription_revenue_by_year(?, ?)", userID, year)
	if err != nil {
		return nil, nil, err
	}

	return coursePurchase, subscription, nil
}

func (r *revenueRepository) GetAdminRevenue(ctx context.Context) (*RevenueBreakdownCents, *RevenueBreakdownCents, error) {
	coursePurchase, err := r.getBreakdown(ctx, "SELECT total_amount, instructor_gross, platform_gross, stripe_fee, instructor_net, platform_net FROM get_admin_course_purchase_revenue()")
	if err != nil {
		return nil, nil, err
	}

	subscription, err := r.getBreakdown(ctx, "SELECT total_amount, instructor_gross, platform_gross, stripe_fee, instructor_net, platform_net FROM get_admin_subscription_revenue()")
	if err != nil {
		return nil, nil, err
	}

	return coursePurchase, subscription, nil
}

func (r *revenueRepository) GetAdminRevenueByDay(ctx context.Context, date time.Time) (*RevenueBreakdownCents, *RevenueBreakdownCents, error) {
	coursePurchase, err := r.getBreakdown(ctx, "SELECT total_amount, instructor_gross, platform_gross, stripe_fee, instructor_net, platform_net FROM get_admin_course_purchase_revenue_by_day(?)", date)
	if err != nil {
		return nil, nil, err
	}

	subscription, err := r.getBreakdown(ctx, "SELECT total_amount, instructor_gross, platform_gross, stripe_fee, instructor_net, platform_net FROM get_admin_subscription_revenue_by_day(?)", date)
	if err != nil {
		return nil, nil, err
	}

	return coursePurchase, subscription, nil
}

func (r *revenueRepository) GetAdminRevenueByMonth(ctx context.Context, year, month int) (*RevenueBreakdownCents, *RevenueBreakdownCents, error) {
	coursePurchase, err := r.getBreakdown(ctx, "SELECT total_amount, instructor_gross, platform_gross, stripe_fee, instructor_net, platform_net FROM get_admin_course_purchase_revenue_by_month(?, ?)", year, month)
	if err != nil {
		return nil, nil, err
	}

	subscription, err := r.getBreakdown(ctx, "SELECT total_amount, instructor_gross, platform_gross, stripe_fee, instructor_net, platform_net FROM get_admin_subscription_revenue_by_month(?, ?)", year, month)
	if err != nil {
		return nil, nil, err
	}

	return coursePurchase, subscription, nil
}

func (r *revenueRepository) GetAdminRevenueByYear(ctx context.Context, year int) (*RevenueBreakdownCents, *RevenueBreakdownCents, error) {
	coursePurchase, err := r.getBreakdown(ctx, "SELECT total_amount, instructor_gross, platform_gross, stripe_fee, instructor_net, platform_net FROM get_admin_course_purchase_revenue_by_year(?)", year)
	if err != nil {
		return nil, nil, err
	}

	subscription, err := r.getBreakdown(ctx, "SELECT total_amount, instructor_gross, platform_gross, stripe_fee, instructor_net, platform_net FROM get_admin_subscription_revenue_by_year(?)", year)
	if err != nil {
		return nil, nil, err
	}

	return coursePurchase, subscription, nil
}

func (r *revenueRepository) GetAdminTransactions(ctx context.Context, limit, offset int, sortBy, sortOrder string) ([]*RevenueTransactionRow, int64, error) {
	const baseUnion = `
		SELECT
			cp.stripe_checkout_session_id AS transaction_id,
			u.id::text AS user_id,
			u.email AS user_email,
			u.username AS user_username,
			u.name AS user_name,
			u.avatar AS user_avatar,
			'stripe' AS method,
			cp.amount::double precision / 100.0 AS amount,
			'course_purchase' AS type,
			COALESCE(cp.purchased_at, cp.created_at) AS created_at
		FROM course_purchases cp
		JOIN users u ON u.id = cp.user_id
		WHERE cp.deleted_at IS NULL
			AND cp.status = 'paid'

		UNION ALL

		SELECT
			s.stripe_subscription_id AS transaction_id,
			u.id::text AS user_id,
			u.email AS user_email,
			u.username AS user_username,
			u.name AS user_name,
			u.avatar AS user_avatar,
			'stripe' AS method,
			s.plan_price::double precision AS amount,
			'subscription' AS type,
			COALESCE(s.started_at, s.created_at) AS created_at
		FROM subscriptions s
		JOIN users u ON u.id = s.user_id
		WHERE s.deleted_at IS NULL
			AND s.status IN ('active', 'trialing', 'past_due', 'canceled')
			AND s.stripe_subscription_id IS NOT NULL
			AND s.stripe_subscription_id <> ''
	`

	countQuery := "SELECT COUNT(1) FROM (" + baseUnion + ") AS transactions"
	var total int64
	if err := r.db.GetDB().WithContext(ctx).Raw(countQuery).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := fmt.Sprintf("%s %s", sortBy, sortOrder)

	listQuery := `
		SELECT transaction_id, user_id, user_email, user_username, user_name, user_avatar, method, amount, type, created_at
		FROM (` + baseUnion + `) AS transactions
		ORDER BY ` + orderClause + `
		LIMIT ? OFFSET ?
	`

	rows := make([]*RevenueTransactionRow, 0)
	if err := r.db.GetDB().WithContext(ctx).Raw(listQuery, limit, offset).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *revenueRepository) getBreakdown(ctx context.Context, query string, args ...any) (*RevenueBreakdownCents, error) {
	var result RevenueBreakdownCents
	if err := r.db.GetDB().WithContext(ctx).Raw(query, args...).Scan(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
