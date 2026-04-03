package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Token     string    `gorm:"type:text;not null"`
	CreatedAt time.Time
	ExpiredAt time.Time `gorm:"index"`
	IsRevoked bool      `gorm:"default:false"`

	User *User `gorm:"foreignKey:UserID;references:ID"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
