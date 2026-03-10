package repository

import (
	"context"
	"elearning-api/model"
	"errors"

	"gorm.io/gorm"
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
