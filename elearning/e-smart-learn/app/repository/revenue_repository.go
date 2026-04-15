package repository

import (
	"context"
	"fmt"
	"strings"
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

	// GetAllTeachersRevenue returns paginated revenue stats for every instructor.
	// startDate / endDate are optional — pass nil to query all-time.
	GetAllTeachersRevenue(
		ctx context.Context,
		startDate *time.Time,
		endDate *time.Time,
		limit int,
		offset int,
		sortOrder string,
	) ([]*TeacherRevenueRow, int64, error)
}

// TeacherRevenueRow is the raw database scan target for the all-teachers query.
// Monetary values are in cents (int64) — the service layer converts to dollars.
type TeacherRevenueRow struct {
	TeacherID     string `gorm:"column:teacher_id"`
	TeacherName   string `gorm:"column:teacher_name"`
	TeacherEmail  string `gorm:"column:teacher_email"`
	TeacherAvatar string `gorm:"column:teacher_avatar"`
	TotalAmount   int64  `gorm:"column:total_amount"`
	StripeFee     int64  `gorm:"column:stripe_fee"`
	InstructorNet int64  `gorm:"column:instructor_net"`
	TotalCourses  int64  `gorm:"column:total_courses"`
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
			cp.total_amount::double precision / 100.0 AS amount,
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

// GetAllTeachersRevenue queries course-purchase revenue aggregated per instructor.
// startDate and endDate are optional — passing nil means no date filter (all-time).
// Results are sorted by total_amount DESC or ASC and paginated.
func (r *revenueRepository) GetAllTeachersRevenue(
	ctx context.Context,
	startDate *time.Time,
	endDate *time.Time,
	limit int,
	offset int,
	sortOrder string,
) ([]*TeacherRevenueRow, int64, error) {
	// Build optional WHERE clauses for date filtering on revenue ledgers.
	var courseDateConditions []string
	var subscriptionDateConditions []string
	var courseDateArgs []any
	var subscriptionDateArgs []any
	if startDate != nil {
		courseDateConditions = append(courseDateConditions, "cprs.purchased_at >= ?")
		subscriptionDateConditions = append(subscriptionDateConditions, "srs.paid_at >= ?")
		courseDateArgs = append(courseDateArgs, *startDate)
		subscriptionDateArgs = append(subscriptionDateArgs, *startDate)
	}
	if endDate != nil {
		courseDateConditions = append(courseDateConditions, "cprs.purchased_at <= ?")
		subscriptionDateConditions = append(subscriptionDateConditions, "srs.paid_at <= ?")
		courseDateArgs = append(courseDateArgs, *endDate)
		subscriptionDateArgs = append(subscriptionDateArgs, *endDate)
	}

	courseDateFilter := ""
	if len(courseDateConditions) > 0 {
		courseDateFilter = "AND " + strings.Join(courseDateConditions, " AND ")
	}

	subscriptionDateFilter := ""
	if len(subscriptionDateConditions) > 0 {
		subscriptionDateFilter = "AND " + strings.Join(subscriptionDateConditions, " AND ")
	}

	// CTE: aggregate ledger-backed revenue per instructor from both
	// course purchases and subscriptions.
	baseCTE := `
		WITH course_revenue AS (
			SELECT
				cprs.instructor_user_id AS teacher_id,
				COALESCE(SUM(cprs.detail_net_amount), 0)::BIGINT AS total_amount,
				COALESCE(SUM(cprs.allocated_stripe_fee), 0)::BIGINT AS stripe_fee,
				COALESCE(SUM(cprs.instructor_net), 0)::BIGINT AS instructor_net
			FROM course_purchase_revenue_shares cprs
			WHERE cprs.deleted_at IS NULL
			  AND cprs.purchased_at IS NOT NULL
			  ` + courseDateFilter + `
			GROUP BY cprs.instructor_user_id
		),
		subscription_revenue AS (
			SELECT
				srs.instructor_user_id AS teacher_id,
				COALESCE(SUM(srs.allocated_amount), 0)::BIGINT AS total_amount,
				COALESCE(SUM(srs.allocated_stripe_fee), 0)::BIGINT AS stripe_fee,
				COALESCE(SUM(srs.instructor_net), 0)::BIGINT AS instructor_net
			FROM subscription_revenue_shares srs
			WHERE srs.deleted_at IS NULL
			  AND srs.paid_at IS NOT NULL
			  ` + subscriptionDateFilter + `
			GROUP BY srs.instructor_user_id
		),
		total_revenue AS (
			SELECT
				teacher_id,
				COALESCE(SUM(total_amount), 0)::BIGINT AS total_amount,
				COALESCE(SUM(stripe_fee), 0)::BIGINT AS stripe_fee,
				COALESCE(SUM(instructor_net), 0)::BIGINT AS instructor_net
			FROM (
				SELECT teacher_id, total_amount, stripe_fee, instructor_net FROM course_revenue
				UNION ALL
				SELECT teacher_id, total_amount, stripe_fee, instructor_net FROM subscription_revenue
			) merged
			GROUP BY teacher_id
		)
	`

	countArgs := make([]any, 0, len(courseDateArgs)+len(subscriptionDateArgs))
	countArgs = append(countArgs, courseDateArgs...)
	countArgs = append(countArgs, subscriptionDateArgs...)

	countQuery := baseCTE + `
		SELECT COUNT(DISTINCT u.id)
		FROM users u
		JOIN roles ro ON ro.id = u.role_id AND ro.deleted_at IS NULL AND ro.name = 'instructor'
		WHERE u.deleted_at IS NULL
	`

	var total int64
	if err := r.db.GetDB().WithContext(ctx).Raw(countQuery, countArgs...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// Allowed sort orders; default to DESC to avoid SQL injection.
	if sortOrder != "asc" {
		sortOrder = "desc"
	}
	orderClause := fmt.Sprintf("total_amount %s", sortOrder)

	listArgs := make([]any, 0, len(courseDateArgs)+len(subscriptionDateArgs)+2)
	listArgs = append(listArgs, courseDateArgs...)
	listArgs = append(listArgs, subscriptionDateArgs...)
	listArgs = append(listArgs, limit, offset)
	listQuery := baseCTE + `
		SELECT
			u.id::text                AS teacher_id,
			u.name                    AS teacher_name,
			u.email                   AS teacher_email,
			COALESCE(u.avatar, '')   AS teacher_avatar,
			COALESCE(tr.total_amount,   0) AS total_amount,
			COALESCE(tr.stripe_fee,     0) AS stripe_fee,
			COALESCE(tr.instructor_net, 0) AS instructor_net,
			COUNT(DISTINCT c.id)          AS total_courses
		FROM users u
		JOIN roles ro
			ON  ro.id         = u.role_id
			AND ro.deleted_at IS NULL
			AND ro.name       = 'instructor'
		LEFT JOIN total_revenue tr ON tr.teacher_id = u.id
		LEFT JOIN courses c
			ON  c.user_id    = u.id
			AND c.deleted_at IS NULL
		WHERE u.deleted_at IS NULL
		GROUP BY u.id, u.name, u.email, u.avatar,
		         tr.total_amount, tr.stripe_fee, tr.instructor_net
		ORDER BY ` + orderClause + `
		LIMIT ? OFFSET ?
	`

	rows := make([]*TeacherRevenueRow, 0)
	if err := r.db.GetDB().WithContext(ctx).Raw(listQuery, listArgs...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}
