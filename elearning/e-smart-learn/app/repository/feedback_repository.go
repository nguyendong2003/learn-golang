package repository

import (
	"context"

	"elearning-api/model"

	"github.com/google/uuid"
)

type FeedbackRepository interface {
	Repository[model.Feedback]
	CountByUserIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]int64, error)
}

type feedbackRepository struct {
	*repository[model.Feedback]
}

func NewFeedbackRepository(db DbRepository) FeedbackRepository {
	return &feedbackRepository{
		repository: NewBaseRepository[model.Feedback](db),
	}
}

func (r *feedbackRepository) CountByUserIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]int64, error) {
	result := make(map[uuid.UUID]int64)
	if len(userIDs) == 0 {
		return result, nil
	}

	type row struct {
		UserID uuid.UUID
		Total  int64
	}

	var rows []row
	err := r.baseQuery(ctx).
		Select("user_id, COUNT(*) as total").
		Where("user_id IN ?", userIDs).
		Group("user_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, it := range rows {
		result[it.UserID] = it.Total
	}
	return result, nil
}
