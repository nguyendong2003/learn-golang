package worker

import (
	"context"
	"sync/atomic"
	"time"

	"elearning-api/pkg"
	"elearning-api/repository"
	"elearning-api/util"
)

type PresignedUploadCleanupWorker struct {
	trackingRepo    repository.PresignedUploadTrackingRepository
	storageProvider pkg.StorageProvider
	running         atomic.Bool
}

func NewPresignedUploadCleanupWorker(
	trackingRepo repository.PresignedUploadTrackingRepository,
	storageProvider pkg.StorageProvider,
) *PresignedUploadCleanupWorker {
	return &PresignedUploadCleanupWorker{
		trackingRepo:    trackingRepo,
		storageProvider: storageProvider,
	}
}

func (w *PresignedUploadCleanupWorker) Start(ctx context.Context) {
	logger := util.WithLayer(util.LayerWorker)
	logger.Info("Starting PresignedUploadCleanupWorker")
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	w.safeRun(ctx)
	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping PresignedUploadCleanupWorker")
			return

		case <-ticker.C:
			w.safeRun(ctx)
		}
	}
}

func (w *PresignedUploadCleanupWorker) safeRun(ctx context.Context) {
	logger := util.WithLayer(util.LayerWorker)
	if !w.running.CompareAndSwap(false, true) {
		logger.Warn("cleanup skipped: previous run still in progress")
		return
	}
	defer w.running.Store(false)
	defer func() {
		if r := recover(); r != nil {
			logger.Error("cleanup panic", "recover", r)
		}
	}()

	jobCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	w.runCleanup(jobCtx)
}

func (w *PresignedUploadCleanupWorker) runCleanup(ctx context.Context) {
	logger := util.WithLayer(util.LayerWorker)
	before := time.Now().Add(-2 * time.Hour)
	total := 0
	for {
		select {
		case <-ctx.Done():
			logger.Warn("cleanup stopped by context")
			return
		default:
		}
		records, err := w.trackingRepo.FindExpiredPending(ctx, before, 100)
		if err != nil {
			logger.Error("Failed to find expired pending uploads", "error", err)
			return
		}
		if len(records) == 0 {
			break
		}
		for _, record := range records {
			if err := w.storageProvider.Delete(ctx, record.ObjectURL); err != nil {
				logger.Warn("Failed to delete orphaned object", "object", record.ObjectURL, "error", err)
				continue
			}
			if err := w.trackingRepo.Delete(ctx, record.ID); err != nil {
				logger.Warn("Failed to delete tracking record", "id", record.ID, "error", err)
				continue
			}
			total++
		}
	}
	logger.Info("cleanup completed", "deleted", total)
}
