package model

import (
	"time"

	"github.com/google/uuid"
)

type CourseEvent struct {
	BaseModel
	CourseID                uuid.UUID `gorm:"type:uuid;not null;index"`
	Name                    string    `gorm:"type:varchar(255);not null"`
	Description             string    `gorm:"type:text"`
	MeetingURL              string    `gorm:"type:text"`
	StartTime               time.Time `gorm:"not null;index"`
	EndTime                 time.Time `gorm:"not null"`
	NotificationBeforeStart int       `gorm:"default:15"` // in minutes

	Course *Course `gorm:"foreignKey:CourseID;references:ID"`
}

func (CourseEvent) TableName() string {
	return "course_events"
}
