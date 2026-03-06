package repository

import (
	"context"
	"elearning-api/model"
)

type UserRepository interface {
	Repository[model.User]

	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmailOrUsername(ctx context.Context, email, username string) (*model.User, error)
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
	return r.Find(ctx, "email = ?", email)
}

func (r *userRepository) GetByUsername(
	ctx context.Context,
	username string,
) (*model.User, error) {
	return r.Find(ctx, "username = ?", username)
}

func (r *userRepository) GetByEmailOrUsername(
	ctx context.Context,
	email string,
	username string,
) (*model.User, error) {
	return r.Find(ctx, "email = ? OR username = ?", email, username)
}
