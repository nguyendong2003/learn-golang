-- +goose Up
-- +goose StatementBegin
-- One payment row per Stripe invoice; prevents duplicate inserts when invoice.paid and
-- subscription sync both call syncPayment (or concurrent webhooks).
CREATE UNIQUE INDEX IF NOT EXISTS idx_payments_stripe_invoice_id_unique
  ON payments (stripe_invoice_id)
  WHERE deleted_at IS NULL
    AND stripe_invoice_id IS NOT NULL
    AND btrim(stripe_invoice_id) <> '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_payments_stripe_invoice_id_unique;
-- +goose StatementEnd
