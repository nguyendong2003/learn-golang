package consts

type BillingCycle string

type SubscriptionStatus string

type PaymentStatus string

type CoursePurchaseStatus string

const (
	BillingCycleMonthly BillingCycle = "monthly"
	BillingCycleYearly  BillingCycle = "yearly"
)

const (
	SubscriptionStatusIncomplete        SubscriptionStatus = "incomplete"
	SubscriptionStatusTrialing          SubscriptionStatus = "trialing"
	SubscriptionStatusActive            SubscriptionStatus = "active"
	SubscriptionStatusPastDue           SubscriptionStatus = "past_due"
	SubscriptionStatusCanceled          SubscriptionStatus = "canceled"
	SubscriptionStatusUnpaid            SubscriptionStatus = "unpaid"
	SubscriptionStatusIncompleteExpired SubscriptionStatus = "incomplete_expired"
)

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

const (
	CoursePurchaseStatusPending  CoursePurchaseStatus = "pending"
	CoursePurchaseStatusPaid     CoursePurchaseStatus = "paid"
	CoursePurchaseStatusFailed   CoursePurchaseStatus = "failed"
	CoursePurchaseStatusRefunded CoursePurchaseStatus = "refunded"
)
