package model

import (
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	UserID    uuid.UUID `gorm:"type:uuid;not null;primaryKey;index"`
	CourseID  uuid.UUID `gorm:"type:uuid;not null;primaryKey;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User   *User   `gorm:"foreignKey:UserID;references:ID"`
	Course *Course `gorm:"foreignKey:CourseID;references:ID"`
}

func (CartItem) TableName() string {
	return "carts"
}
