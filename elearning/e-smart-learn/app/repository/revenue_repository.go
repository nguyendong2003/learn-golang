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
	// Build optional WHERE clauses for date filtering on course_purchases.
	var dateConditions []string
	var dateArgs []any
	if startDate != nil {
		dateConditions = append(dateConditions, "cp.purchased_at >= ?")
		dateArgs = append(dateArgs, *startDate)
	}
	if endDate != nil {
		dateConditions = append(dateConditions, "cp.purchased_at <= ?")
		dateArgs = append(dateArgs, *endDate)
	}

	dateFilter := ""
	if len(dateConditions) > 0 {
		dateFilter = "AND " + strings.Join(dateConditions, " AND ")
	}

	// CTE: aggregate course-purchase revenue per instructor.
	// stripe_fee is prorated per course-detail line proportionally to its price.
	baseCTE := `
		WITH course_revenue AS (
			SELECT
				c.user_id AS teacher_id,
				SUM(cpd.price)                                                         AS total_amount,
				SUM(cp.stripe_fee * cpd.price / NULLIF(cp.amount, 0))                  AS stripe_fee,
				SUM(cpd.price) - SUM(cp.stripe_fee * cpd.price / NULLIF(cp.amount, 0)) AS instructor_net
			FROM course_purchase_details cpd
			JOIN course_purchases cp
				ON  cp.id         = cpd.course_purchase_id
				AND cp.status     = 'paid'
				AND cp.deleted_at IS NULL
				` + dateFilter + `
			JOIN courses c
				ON  c.id         = cpd.course_id
				AND c.deleted_at IS NULL
			WHERE cpd.deleted_at IS NULL
			GROUP BY c.user_id
		)
	`

	// Prepare args: date args first, then will add limit/offset for list query.
	countArgs := make([]any, len(dateArgs))
	copy(countArgs, dateArgs)

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

	listArgs := append(dateArgs, limit, offset) //nolint:gocritic
	listQuery := baseCTE + `
		SELECT
			u.id::text                AS teacher_id,
			u.name                    AS teacher_name,
			u.email                   AS teacher_email,
			COALESCE(u.avatar, '')   AS teacher_avatar,
			COALESCE(cr.total_amount,   0) AS total_amount,
			COALESCE(cr.stripe_fee,     0) AS stripe_fee,
			COALESCE(cr.instructor_net, 0) AS instructor_net,
			COUNT(DISTINCT c.id)          AS total_courses
		FROM users u
		JOIN roles ro
			ON  ro.id         = u.role_id
			AND ro.deleted_at IS NULL
			AND ro.name       = 'instructor'
		LEFT JOIN course_revenue cr ON cr.teacher_id = u.id
		LEFT JOIN courses c
			ON  c.user_id    = u.id
			AND c.deleted_at IS NULL
		WHERE u.deleted_at IS NULL
		GROUP BY u.id, u.name, u.email, u.avatar,
		         cr.total_amount, cr.stripe_fee, cr.instructor_net
		ORDER BY ` + orderClause + `
		LIMIT ? OFFSET ?
	`

	rows := make([]*TeacherRevenueRow, 0)
	if err := r.db.GetDB().WithContext(ctx).Raw(listQuery, listArgs...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}
