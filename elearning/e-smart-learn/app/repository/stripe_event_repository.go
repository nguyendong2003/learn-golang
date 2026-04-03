package repository

import (
	"context"
	"elearning-api/model"
)

type StripeEventRepository interface {
	Repository[model.StripeEvent]
	GetByEventID(ctx context.Context, eventID string) (*model.StripeEvent, error)
}

type stripeEventRepository struct {
	*repository[model.StripeEvent]
}

func NewStripeEventRepository(db DbRepository) StripeEventRepository {
	return &stripeEventRepository{
		repository: NewBaseRepository[model.StripeEvent](db),
	}
}

func (r *stripeEventRepository) GetByEventID(ctx context.Context, eventID string) (*model.StripeEvent, error) {
	return r.Find(ctx, "event_id = ?", nil, eventID)
}
