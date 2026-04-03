package repository

import (
	"context"
	"elearning-api/model"
)

type PlanRepository interface {
	Repository[model.Plan]
	ListActivePlans(ctx context.Context) ([]*model.Plan, error)
}

type planRepository struct {
	*repository[model.Plan]
}

func NewPlanRepository(db DbRepository) PlanRepository {
	return &planRepository{
		repository: NewBaseRepository[model.Plan](db),
	}
}

func (r *planRepository) ListActivePlans(ctx context.Context) ([]*model.Plan, error) {
	return r.FindAll(ctx, "is_active = ?", nil, true)
}
