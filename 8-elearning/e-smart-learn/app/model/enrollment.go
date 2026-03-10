package model

import (
	"time"

	"github.com/google/uuid"
)

type Enrollment struct {
	UserID      uuid.UUID `gorm:"type:uuid;not null;primaryKey;index"`
	CourseID    uuid.UUID `gorm:"type:uuid;not null;primaryKey;index"`
	EnrolledAt  time.Time `gorm:"autoCreateTime"`
	CompletedAt *time.Time
	DeletedAt   *time.Time `gorm:"index"`

	User   *User   `gorm:"foreignKey:UserID;references:ID"`
	Course *Course `gorm:"foreignKey:CourseID;references:ID"`
}

func (Enrollment) TableName() string {
	return "enrollments"
}

type Feedback struct {
	BaseModel
	UserID   uuid.UUID `gorm:"type:uuid;not null;index"`
	CourseID uuid.UUID `gorm:"type:uuid;not null;index"`
	Rate     int       `gorm:"type:int;default:5"` 
	Content  string    `gorm:"type:text"`

	User   *User   `gorm:"foreignKey:UserID;references:ID"`
	Course *Course `gorm:"foreignKey:CourseID;references:ID"`
}

func (Feedback) TableName() string {
	return "feedbacks"
}
