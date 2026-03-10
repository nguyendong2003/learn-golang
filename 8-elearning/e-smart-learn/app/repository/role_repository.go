package repository

import (
	"context"
	"elearning-api/model"
)

type RoleRepository interface {
	Repository[model.Role]

	GetAll(ctx context.Context) ([]*model.Role, error)
	GetAllWithPermissions(ctx context.Context) ([]*model.Role, error)
	GetByName(ctx context.Context, name string) (*model.Role, error)
	GetByNameWithPermissions(ctx context.Context, name string) (*model.Role, error)
}

type roleRepository struct {
	*repository[model.Role]
}

func NewRoleRepository(db DbRepository) RoleRepository {
	return &roleRepository{
		repository: NewBaseRepository[model.Role](db),
	}
}

func (r *roleRepository) GetAll(ctx context.Context) ([]*model.Role, error) {
	return r.FindAll(ctx, "")
}

func (r *roleRepository) GetAllWithPermissions(ctx context.Context) ([]*model.Role, error) {
	var roles []*model.Role
	if err := r.baseQuery(ctx).
		Preload("Permissions").
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*model.Role, error) {
	return r.Find(ctx, "name = ?", name)
}

func (r *roleRepository) GetByNameWithPermissions(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	if err := r.baseQuery(ctx).
		Preload("Permissions").
		Where("name = ?", name).
		First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
