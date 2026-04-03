package repository

import (
	"context"
	"errors"

	"elearning-api/consts"
	"elearning-api/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BlogRepository interface {
	Repository[model.Blog]
	FindBySlug(ctx context.Context, slug string, preloads []Preload) (*model.Blog, error)
	IncrementView(ctx context.Context, id uuid.UUID) error
	GetStats(ctx context.Context) (*model.BlogStats, error)
}

type blogRepository struct {
	*repository[model.Blog]
}

func NewBlogRepository(db DbRepository) BlogRepository {
	return &blogRepository{
		repository: NewBaseRepository[model.Blog](db),
	}
}

func (r *blogRepository) FindBySlug(ctx context.Context, slug string, preloads []Preload) (*model.Blog, error) {
	var blog model.Blog
	query := r.baseQuery(ctx).Where("slug = ? AND status = ?", slug, consts.BlogStatusPublished)

	if len(preloads) > 0 {
		applyPreloads(query, preloads)
	}

	if err := query.First(&blog).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *blogRepository) IncrementView(ctx context.Context, id uuid.UUID) error {
	return r.db.GetDB().WithContext(ctx).
		Model(&model.Blog{}).
		Where("id = ?", id).
		Update("view_total", gorm.Expr("view_total + ?", 1)).Error
}

func (r *blogRepository) GetStats(ctx context.Context) (*model.BlogStats, error) {
	var stats model.BlogStats
	err := r.baseQuery(ctx).
		Select(`
			COUNT(*) AS total_articles,
			COALESCE(SUM(view_total), 0) AS total_views,
			COALESCE(SUM(CASE WHEN status = ? THEN 1 ELSE 0 END), 0) AS published,
			COALESCE(SUM(CASE WHEN status = ? THEN 1 ELSE 0 END), 0) AS drafts,
			COALESCE(SUM(CASE WHEN status = ? THEN 1 ELSE 0 END), 0) AS scheduled,
			COALESCE(SUM(CASE WHEN view_total > 0 THEN 1 ELSE 0 END), 0) AS engaged
		`,
			consts.BlogStatusPublished,
			consts.BlogStatusDraft,
			consts.BlogStatusScheduled,
		).
		Where("deleted_at IS NULL").
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
