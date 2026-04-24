package repository

import (
	"context"
	"time"

	"elearning-api/consts"
	"elearning-api/model"

	"github.com/google/uuid"
)

type EnrollmentRepository interface {
	Repository[model.Enrollment]
	AddLearnCourse(ctx context.Context, userID uuid.UUID, courseID uuid.UUID, lessonID uuid.UUID) (*model.Enrollment, error)
	CheckEnrollment(ctx context.Context, userId, courseId uuid.UUID) (bool, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Enrollment, int64, error)
	GetCourseCompletionTotals(ctx context.Context) (int64, int64, error)
	EnrollIfNotExists(ctx context.Context, userID, courseID uuid.UUID, enrollmentType consts.EnrollmentType) (*model.Enrollment, error)
	CancelActiveSubscriptionEnrollmentsByUser(ctx context.Context, userID uuid.UUID) error
	MarkCourseCompleted(ctx context.Context, userID, courseID uuid.UUID) error
	GetCourseEnrollmentCount(ctx context.Context, userID uuid.UUID) (int64, int64, error)
	GetTopCategoryByUser(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	GetEnrolledCategoryIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

type enrollmentRepository struct {
	*repository[model.Enrollment]
}

func NewEnrollmentRepository(db DbRepository) EnrollmentRepository {
	return &enrollmentRepository{
		repository: NewBaseRepository[model.Enrollment](db),
	}
}

func (r *enrollmentRepository) AddLearnCourse(ctx context.Context, userID uuid.UUID, courseID uuid.UUID, lessonID uuid.UUID) (*model.Enrollment, error) {
	enrollment, err := r.Find(ctx, "user_id=? AND course_id=?", nil, userID, courseID)
	if err != nil {
		return nil, err
	}
	enrollment.LearnedLessonIds = append(enrollment.LearnedLessonIds, lessonID)
	if err := r.db.GetDB().WithContext(ctx).Save(enrollment).Error; err != nil {
		return nil, err
	}
	return enrollment, nil
}

func (r *enrollmentRepository) CheckEnrollment(ctx context.Context, userId, courseId uuid.UUID) (bool, error) {
	var count int64
	if err := r.baseQuery(ctx).
		Where("user_id = ? AND course_id = ? AND canceled_at IS NULL", userId, courseId).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *enrollmentRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Enrollment, int64, error) {
	preloads := []Preload{
		Course,
		PreloadPath(Course, Chapters, Lessons),
	}
	return r.List(ctx, limit, offset, "enrolled_at desc", "user_id = ? AND canceled_at IS NULL", preloads, userID)
}

func (r *enrollmentRepository) GetCourseCompletionTotals(ctx context.Context) (int64, int64, error) {
	type completionTotals struct {
		TotalCompletedLessons int64 `gorm:"column:total_completed_lessons"`
		TotalLessons          int64 `gorm:"column:total_lessons"`
	}

	var totals completionTotals
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
		)
		SELECT
			COALESCE(SUM(LEAST(COALESCE(array_length(e.learned_lesson_ids, 1), 0), COALESCE(lpc.total_lessons, 0))), 0) AS total_completed_lessons,
			COALESCE(SUM(COALESCE(lpc.total_lessons, 0)), 0) AS total_lessons
		FROM enrollments e
		LEFT JOIN lessons_per_course lpc ON lpc.course_id = e.course_id
		WHERE e.deleted_at IS NULL
	`

	if err := r.db.GetDB().WithContext(ctx).Raw(query).Scan(&totals).Error; err != nil {
		return 0, 0, err
	}

	return totals.TotalCompletedLessons, totals.TotalLessons, nil
}

func (r *enrollmentRepository) GetCourseEnrollmentCount(ctx context.Context, userID uuid.UUID) (int64, int64, error) {
	type enrollmentCount struct {
		TotalEnrollments int64 `gorm:"column:total_enrollments"`
		TotalCompletions int64 `gorm:"column:total_completions"`
	}

	var count enrollmentCount
	query := `
		SELECT
			COUNT(*) AS total_enrollments,
			COUNT(CASE WHEN is_completed THEN 1 END) AS total_completions
		FROM enrollments
		WHERE user_id = ? AND deleted_at IS NULL
	`

	if err := r.db.GetDB().WithContext(ctx).Raw(query, userID).Scan(&count).Error; err != nil {
		return 0, 0, err
	}

	return count.TotalEnrollments, count.TotalCompletions, nil
}

func (r *enrollmentRepository) EnrollIfNotExists(ctx context.Context, userID, courseID uuid.UUID, enrollmentType consts.EnrollmentType) (*model.Enrollment, error) {
	enrollment, err := r.Find(ctx, "user_id = ? AND course_id = ? AND canceled_at IS NULL", nil, userID, courseID)
	if err != nil {
		return nil, err
	}
	if enrollment != nil {
		return enrollment, nil
	}

	enrollment = &model.Enrollment{
		UserID:     userID,
		CourseID:   courseID,
		EnrolledAt: time.Now().UTC(),
		Type:       string(enrollmentType),
	}

	if err := r.db.GetDB().WithContext(ctx).Create(enrollment).Error; err != nil {
		return nil, err
	}

	return enrollment, nil
}

func (r *enrollmentRepository) CancelActiveSubscriptionEnrollmentsByUser(ctx context.Context, userID uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.GetDB().WithContext(ctx).
		Model(&model.Enrollment{}).
		Where("user_id = ? AND type = ? AND canceled_at IS NULL", userID, consts.EnrollmentTypeSubscription).
		Update("canceled_at", &now).Error
}

func (r *enrollmentRepository) MarkCourseCompleted(ctx context.Context, userID, courseID uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.GetDB().WithContext(ctx).
		Model(&model.Enrollment{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Updates(map[string]any{
			"is_completed": true,
			"completed_at": &now,
		}).Error
}

func (r *enrollmentRepository) GetTopCategoryByUser(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	type catRow struct {
		CategoryID uuid.UUID `gorm:"column:category_id"`
		Total      int64     `gorm:"column:total"`
	}

	var row catRow
	query := r.db.GetDB().WithContext(ctx).
		Table("enrollments e").
		Select("c.category_id as category_id, COUNT(*) as total").
		Joins("JOIN courses c ON c.id = e.course_id").
		Where("e.user_id = ? AND e.canceled_at IS NULL AND c.deleted_at IS NULL", userID).
		Group("c.category_id").
		Order("total DESC").
		Limit(1)

	if err := query.Scan(&row).Error; err != nil {
		return uuid.Nil, err
	}

	if row.CategoryID == uuid.Nil {
		return uuid.Nil, nil
	}

	return row.CategoryID, nil
}

// Return list of category IDs the user has enrolled in, ordered by most enrolled to least enrolled
func (r *enrollmentRepository) GetEnrolledCategoryIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	type catRow struct {
		CategoryID uuid.UUID `gorm:"column:category_id"`
		Total      int64     `gorm:"column:total"`
	}

	var rows []catRow
	query := r.db.GetDB().WithContext(ctx).
		Table("enrollments e").
		Select("DISTINCT c.category_id as category_id, COUNT(*) as total").
		Joins("JOIN courses c ON c.id = e.course_id").
		Where("e.user_id = ? AND e.canceled_at IS NULL AND c.deleted_at IS NULL", userID).
		Group("c.category_id").
		Order("total DESC")

	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}

	res := make([]uuid.UUID, 0, len(rows))
	for _, rrow := range rows {
		if rrow.CategoryID != uuid.Nil {
			res = append(res, rrow.CategoryID)
		}
	}

	return res, nil
}
