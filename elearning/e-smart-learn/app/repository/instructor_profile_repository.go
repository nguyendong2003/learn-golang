package repository

import (
	"context"
	"elearning-api/model"
	"time"
)

type InstructorProfileRepository interface {
	Repository[model.InstructorProfile]
	GetGrowthStatistics(ctx context.Context, since, until time.Time) (*TeacherGrowthStatisticsRow, error)
}

type instructorProfileRepository struct {
	*repository[model.InstructorProfile]
}

type TeacherGrowthStatisticsRow struct {
	TotalVerifiedTeachers int64  `gorm:"column:total_verified_teachers"`
	NewThisQuarter        int64  `gorm:"column:new_this_quarter"`
	TopCategoryName       string `gorm:"column:top_category_name"`
	TopCategoryCount      int64  `gorm:"column:top_category_count"`
}

func NewInstructorProfileRepository(db DbRepository) InstructorProfileRepository {
	return &instructorProfileRepository{
		repository: NewBaseRepository[model.InstructorProfile](db),
	}
}

func (r *instructorProfileRepository) GetGrowthStatistics(ctx context.Context, since, until time.Time) (*TeacherGrowthStatisticsRow, error) {
	var stats TeacherGrowthStatisticsRow

	query := `
		WITH approved_profiles AS (
			SELECT
				ip.id,
				ip.category_id,
				ip.created_at
			FROM instructor_profiles ip
			WHERE ip.deleted_at IS NULL
			  AND ip.status = 'approved'
		),
		category_counts AS (
			SELECT
				c.name AS category_name,
				COUNT(*)::BIGINT AS category_count
			FROM approved_profiles ap
			JOIN categories c ON c.id = ap.category_id AND c.deleted_at IS NULL
			GROUP BY c.name
		)
		SELECT
			COALESCE(COUNT(*), 0)::BIGINT AS total_verified_teachers,
			COALESCE(COUNT(*) FILTER (WHERE created_at >= ? AND created_at < ?), 0)::BIGINT AS new_this_quarter,
			COALESCE((
				SELECT category_name
				FROM category_counts
				ORDER BY category_count DESC, category_name ASC
				LIMIT 1
			), '') AS top_category_name,
			COALESCE((
				SELECT category_count
				FROM category_counts
				ORDER BY category_count DESC, category_name ASC
				LIMIT 1
			), 0)::BIGINT AS top_category_count
		FROM approved_profiles;
	`

	if err := r.db.GetDB().WithContext(ctx).Raw(query, since, until).Scan(&stats).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
