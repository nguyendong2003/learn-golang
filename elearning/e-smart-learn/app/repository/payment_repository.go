package repository

import (
	"context"
	"elearning-api/model"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Repository[model.Payment]
	GetByStripeInvoiceID(ctx context.Context, stripeInvoiceID string, preloads []Preload) (*model.Payment, error)
	GetSucceededPaymentBySubscriptionAndTime(ctx context.Context, subscriptionID uuid.UUID, asOf time.Time, preloads []Preload) (*model.Payment, error)
}

type paymentRepository struct {
	*repository[model.Payment]
}

func NewPaymentRepository(db DbRepository) PaymentRepository {
	return &paymentRepository{
		repository: NewBaseRepository[model.Payment](db),
	}
}

func (r *paymentRepository) GetByStripeInvoiceID(ctx context.Context, stripeInvoiceID string, preloads []Preload) (*model.Payment, error) {
	return r.Find(ctx, "stripe_invoice_id = ?", preloads, stripeInvoiceID)
}

func (r *paymentRepository) GetSucceededPaymentBySubscriptionAndTime(ctx context.Context, subscriptionID uuid.UUID, asOf time.Time, preloads []Preload) (*model.Payment, error) {
	query := `
		subscription_id = ?
		AND status = 'succeeded'
		AND paid_at IS NOT NULL
		AND COALESCE(billing_period_start, paid_at) <= ?
		AND COALESCE(billing_period_end, paid_at + interval '1 second') > ?
	`
	db := r.baseQuery(ctx).Where(query, subscriptionID, asOf, asOf).Order("paid_at DESC")
	if len(preloads) > 0 {
		db = applyPreloads(db, preloads)
	}

	var payment model.Payment
	if err := db.First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}
