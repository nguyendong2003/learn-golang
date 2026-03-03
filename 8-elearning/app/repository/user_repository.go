package repository

import (
	"context"
	"elearning-api/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Repository[model.User]
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}

type userRepository struct {
	*repository[model.User]
}

func NewUserRepository(db DbRepository) UserRepository {
	return &userRepository{
		repository: NewBaseRepository[model.User](db),
	}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	if err := r.db.GetDB().WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	if err := r.db.GetDB().WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
