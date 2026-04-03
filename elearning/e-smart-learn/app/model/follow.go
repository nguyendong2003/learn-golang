package model

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	FollowerID uuid.UUID `gorm:"type:uuid;not null;primaryKey;index"`
	FolloweeID uuid.UUID `gorm:"type:uuid;not null;primaryKey;index"`
	CreatedAt  time.Time

	Follower *User `gorm:"foreignKey:FollowerID;references:ID"`
	Followee *User `gorm:"foreignKey:FolloweeID;references:ID"`
}

func (Follow) TableName() string {
	return "follows"
}
