package repository

import (
	"context"
	"elearning-api/model"
)

type PaymentRepository interface {
	Repository[model.Payment]
	GetByStripeInvoiceID(ctx context.Context, stripeInvoiceID string, preloads []Preload) (*model.Payment, error)
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
