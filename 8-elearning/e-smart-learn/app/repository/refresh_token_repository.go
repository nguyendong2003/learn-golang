package repository

import (
	"context"
	"elearning-api/model"
	"errors"

	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Repository[model.RefreshToken]

	GetByToken(ctx context.Context, token string) (*model.RefreshToken, error)
}

type refreshTokenRepository struct {
	*repository[model.RefreshToken]
}

func NewRefreshTokenRepository(db DbRepository) RefreshTokenRepository {
	return &refreshTokenRepository{
		repository: NewBaseRepository[model.RefreshToken](db),
	}
}

func (r *refreshTokenRepository) GetByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken
	err := r.baseQuery(ctx).
		Where("token = ?", token).
		Preload("User").
		First(&refreshToken).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}
