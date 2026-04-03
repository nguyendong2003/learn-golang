package model

type StripeEvent struct {
	BaseModel

	EventID   string `gorm:"type:varchar(255);not null;uniqueIndex"`
	EventType string `gorm:"type:varchar(100);not null;index"`
}

func (StripeEvent) TableName() string {
	return "stripe_events"
}
