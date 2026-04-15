package repository

import (
	"context"
	"strings"
	"time"

	"elearning-api/consts"
	"elearning-api/model"

	"github.com/google/uuid"
)

type CourseRepository interface {
	Repository[model.Course]
	GetStatistics(ctx context.Context) (*model.CourseStatistics, error)
	GetInstructorTaughtCourseRevenue(
		ctx context.Context,
		userID uuid.UUID,
		limit int,
		offset int,
		sortBy string,
		sortOrder string,
	) ([]*model.InstructorTaughtCourseRevenue, int64, error)
	CountCreatedSince(ctx context.Context, since time.Time) (int64, error)
}

type courseRepository struct {
	*repository[model.Course]
}

func NewCourseRepository(db DbRepository) CourseRepository {
	return &courseRepository{
		repository: NewBaseRepository[model.Course](db),
	}
}

func (r *courseRepository) GetStatistics(ctx context.Context) (*model.CourseStatistics, error) {
	var stats model.CourseStatistics
	err := r.baseQuery(ctx).
		Select(`
			COUNT(*) AS total_courses,
			COALESCE(SUM(CASE WHEN status = 'pending_review' THEN 1 ELSE 0 END), 0) AS pending_reviews,
			COALESCE(SUM(CASE WHEN status = 'draft' THEN 1 ELSE 0 END), 0) AS drafts,
			COALESCE(SUM(CASE WHEN status = 'published' THEN 1 ELSE 0 END), 0) AS published,
			COALESCE(SUM(CASE WHEN status = 'archived' THEN 1 ELSE 0 END), 0) AS archived
		`).Scan(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *courseRepository) GetInstructorTaughtCourseRevenue(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
	offset int,
	sortBy string,
	sortOrder string,
) ([]*model.InstructorTaughtCourseRevenue, int64, error) {
	var total int64
	if err := r.baseQuery(ctx).
		Where("user_id = ? AND status = ?", userID, consts.CoursePublished).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	allowedSort := map[string]string{
		"created_at":    "c.created_at",
		"title":         "c.title",
		"status":        "c.status",
		"total_student": "c.total_student",
		"revenue":       "revenue",
	}

	sortColumn, ok := allowedSort[sortBy]
	if !ok {
		sortColumn = "c.created_at"
	}

	order := strings.ToUpper(sortOrder)
	if order != "ASC" {
		order = "DESC"
	}

	orderClause := sortColumn + " " + order

	query := r.db.GetDB().WithContext(ctx).
		Table("courses c").
		Select(`
			c.id AS course_id,
			c.title,
			c.slug,
			c.image,
			c.status,
			c.total_student,
			c.created_at,
			COALESCE(SUM(CASE WHEN cp.id IS NOT NULL THEN cpd.price_final ELSE 0 END), 0) AS revenue
		`).
		Joins("LEFT JOIN course_purchase_details cpd ON cpd.course_id = c.id AND cpd.deleted_at IS NULL").
		Joins("LEFT JOIN course_purchases cp ON cp.id = cpd.course_purchase_id AND cp.deleted_at IS NULL AND cp.status = ?", consts.CoursePurchaseStatusPaid).
		Where("c.user_id = ? AND c.deleted_at IS NULL AND c.status = ?", userID, consts.CoursePublished).
		Group("c.id").
		Order(orderClause)

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	var rows []*model.InstructorTaughtCourseRevenue
	if err := query.Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *courseRepository) CountCreatedSince(ctx context.Context, since time.Time) (int64, error) {
	return r.Count(ctx, "created_at >= ? AND deleted_at IS NULL", since)
}
