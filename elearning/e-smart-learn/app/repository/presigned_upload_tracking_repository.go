package repository

import (
	"context"
	"time"

	"elearning-api/consts"
	"elearning-api/model"
)

type PresignedUploadTrackingRepository interface {
	Repository[model.PresignedUploadTracking]
	FindExpiredPending(ctx context.Context, before time.Time, limit int) ([]*model.PresignedUploadTracking, error)
	ConfirmByObjectURL(ctx context.Context, objectURL string) error
}

type presignedUploadTrackingRepository struct {
	*repository[model.PresignedUploadTracking]
}

func NewPresignedUploadTrackingRepository(db DbRepository) PresignedUploadTrackingRepository {
	return &presignedUploadTrackingRepository{
		repository: NewBaseRepository[model.PresignedUploadTracking](db),
	}
}

func (r *presignedUploadTrackingRepository) FindExpiredPending(ctx context.Context, before time.Time, limit int) ([]*model.PresignedUploadTracking, error) {
	var results []*model.PresignedUploadTracking

	err := r.baseQuery(ctx).
		Where("status = ? AND created_at < ?", consts.PresignedUploadStatusPending, before).
		Limit(limit).
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *presignedUploadTrackingRepository) ConfirmByObjectURL(ctx context.Context, objectURL string) error {
	result := r.baseQuery(ctx).
		Where("object_url = ? AND status = ?", objectURL, consts.PresignedUploadStatusPending).
		Update("status", consts.PresignedUploadStatusConfirmed)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
