package repository

import (
	"context"

	"elearning-api/model"

	"github.com/google/uuid"
)

type FollowRepository interface {
	Repository[model.Follow]
	// Check if a follow relationship exists between follower and followee
	Exists(ctx context.Context, followerID, followeeID uuid.UUID) (bool, error)
	DeleteByPair(ctx context.Context, followerID, followeeID uuid.UUID) error
	ListFollowers(ctx context.Context, followeeID uuid.UUID, limit, offset int, order string) ([]*model.Follow, int64, error)
	ListFollowings(ctx context.Context, followerID uuid.UUID, limit, offset int, order string) ([]*model.Follow, int64, error)
}

type followRepository struct {
	*repository[model.Follow]
}

func NewFollowRepository(db DbRepository) FollowRepository {
	return &followRepository{
		repository: NewBaseRepository[model.Follow](db),
	}
}

func (r *followRepository) Exists(ctx context.Context, followerID, followeeID uuid.UUID) (bool, error) {
	return r.CheckExists(ctx, "follower_id = ? AND followee_id = ?", followerID, followeeID)
}

func (r *followRepository) DeleteByPair(ctx context.Context, followerID, followeeID uuid.UUID) error {
	return r.baseQuery(ctx).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Delete(&model.Follow{}).
		Error
}

func (r *followRepository) ListFollowers(ctx context.Context, followeeID uuid.UUID, limit, offset int, order string) ([]*model.Follow, int64, error) {
	return r.List(ctx, limit, offset, order, "followee_id = ?", []Preload{PreloadPath(Follower)}, followeeID)
}

func (r *followRepository) ListFollowings(ctx context.Context, followerID uuid.UUID, limit, offset int, order string) ([]*model.Follow, int64, error) {
	return r.List(ctx, limit, offset, order, "follower_id = ?", []Preload{PreloadPath(Followee)}, followerID)
}
