package repository

import (
	"context"
	"errors"
	"time"

	"elearning-api/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Repository[model.User]

	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmailOrUsername(ctx context.Context, email, username string) (*model.User, error)
	GetByOAuth(ctx context.Context, provider, oauthID string) (*model.User, error)
	GetList(ctx context.Context, limit, offset int) ([]*model.UserDirectoryRow, int64, error)
	CountActiveByRoleName(ctx context.Context, roleName string) (int64, error)
	CountActiveByRoleNameBefore(ctx context.Context, roleName string, before time.Time) (int64, error)
}

type userRepository struct {
	*repository[model.User]
}

func NewUserRepository(db DbRepository) UserRepository {
	return &userRepository{
		repository: NewBaseRepository[model.User](db),
	}
}

func (r *userRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {
	var user model.User
	err := r.baseQuery(ctx).
		Where("email = ?", email).
		Preload("Role").
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByUsername(
	ctx context.Context,
	username string,
) (*model.User, error) {
	var user model.User
	err := r.baseQuery(ctx).
		Where("username = ?", username).
		Preload("Role").
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByEmailOrUsername(
	ctx context.Context,
	email string,
	username string,
) (*model.User, error) {
	var user model.User
	err := r.baseQuery(ctx).
		Where("email = ? OR username = ?", email, username).
		Preload("Role").
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByOAuth(
	ctx context.Context,
	provider string,
	oauthID string,
) (*model.User, error) {
	var user model.User
	err := r.baseQuery(ctx).
		Where("oauth_provider = ? AND oauth_id = ?", provider, oauthID).
		Preload("Role").
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetList(ctx context.Context, limit, offset int) ([]*model.UserDirectoryRow, int64, error) {
	countQuery := `
		WITH instructor_pending AS (
			SELECT
				ip.user_id,
				BOOL_OR(ip.status = 'pending_review') AS has_pending_instructor_application
			FROM instructor_profiles ip
			WHERE ip.deleted_at IS NULL
			GROUP BY ip.user_id
		)
		SELECT COUNT(*)
		FROM users u
		LEFT JOIN roles r ON r.id = u.role_id AND r.deleted_at IS NULL
		LEFT JOIN instructor_pending ip ON ip.user_id = u.id
		WHERE u.deleted_at IS NULL

	`

	var total int64
	if err := r.db.GetDB().WithContext(ctx).Raw(countQuery).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	query := `
		WITH lessons_per_course AS (
			SELECT
				c.id AS course_id,
				COUNT(l.id) AS total_lessons
			FROM courses c
			LEFT JOIN chapters ch ON ch.course_id = c.id AND ch.deleted_at IS NULL
			LEFT JOIN lessons l ON l.chapter_id = ch.id AND l.deleted_at IS NULL
			WHERE c.deleted_at IS NULL
			GROUP BY c.id
		),
		student_progress AS (
			SELECT
				e.user_id,
				COALESCE(SUM(LEAST(COALESCE(array_length(e.learned_lesson_ids, 1), 0), COALESCE(lpc.total_lessons, 0))), 0) AS completed_lessons,
				COALESCE(SUM(COALESCE(lpc.total_lessons, 0)), 0) AS total_lessons
			FROM enrollments e
			LEFT JOIN lessons_per_course lpc ON lpc.course_id = e.course_id
			WHERE e.deleted_at IS NULL
			GROUP BY e.user_id
		),
		instructor_courses AS (
			SELECT
				c.user_id,
				COUNT(1) AS total_courses_taught,
				COUNT(1) FILTER (WHERE c.status = 'published') AS active_courses
			FROM courses c
			WHERE c.deleted_at IS NULL
			GROUP BY c.user_id
		),
		instructor_pending AS (
			SELECT
				ip.user_id,
				BOOL_OR(ip.status = 'pendilng_review') AS has_pending_instructor_application
			FROM instructor_profiles ip
			WHERE ip.deleted_at IS NULL
			GROUP BY ip.user_id
		)
		SELECT
			u.id::text AS user_id,
			u.name,
			u.email,
			COALESCE(u.avatar, '') AS avatar,
			COALESCE(r.name, '') AS role_name,
			u.is_active,
			COALESCE(ip.has_pending_instructor_application, false) AS has_pending_instructor_application,
			COALESCE(ic.active_courses, 0) AS active_courses,
			COALESCE(ic.total_courses_taught, 0) AS total_courses_taught,
			COALESCE(sp.completed_lessons, 0) AS completed_lessons,
			COALESCE(sp.total_lessons, 0) AS total_lessons
		FROM users u
		LEFT JOIN roles r ON r.id = u.role_id AND r.deleted_at IS NULL
		LEFT JOIN student_progress sp ON sp.user_id = u.id
		LEFT JOIN instructor_courses ic ON ic.user_id = u.id
		LEFT JOIN instructor_pending ip ON ip.user_id = u.id
		WHERE u.deleted_at IS NULL
		ORDER BY u.created_at DESC
		LIMIT ? OFFSET ?
	`

	var rows []*model.UserDirectoryRow
	if err := r.db.GetDB().WithContext(ctx).Raw(query, limit, offset).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *userRepository) CountActiveByRoleName(ctx context.Context, roleName string) (int64, error) {
	query := `
		SELECT COUNT(1)
		FROM users u
		JOIN roles r ON r.id = u.role_id AND r.deleted_at IS NULL
		WHERE u.deleted_at IS NULL
		  AND u.is_active = TRUE
		  AND r.name = ?
	`

	var total int64
	if err := r.db.GetDB().WithContext(ctx).Raw(query, roleName).Scan(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *userRepository) CountActiveByRoleNameBefore(ctx context.Context, roleName string, before time.Time) (int64, error) {
	query := `
		SELECT COUNT(1)
		FROM users u
		JOIN roles r ON r.id = u.role_id AND r.deleted_at IS NULL
		WHERE u.deleted_at IS NULL
		  AND u.is_active = TRUE
		  AND r.name = ?
		  AND u.created_at < ?
	`

	var total int64
	if err := r.db.GetDB().WithContext(ctx).Raw(query, roleName, before).Scan(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
