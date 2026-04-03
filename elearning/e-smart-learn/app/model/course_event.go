package model

import (
	"time"

	"elearning-api/repository/dbtypes"

	"github.com/google/uuid"
)

type CourseEvent struct {
	BaseModel
	CourseID                uuid.UUID           `gorm:"type:uuid;not null;index"`
	Name                    string              `gorm:"type:varchar(255);not null"`
	Description             string              `gorm:"type:text"`
	Location                string              `gorm:"type:varchar(255)"`
	RoomToken               string              `gorm:"type:varchar(255);unique;not null"`
	StartTime               time.Time           `gorm:"not null;index"`
	EndTime                 time.Time           `gorm:"not null"`
	NotificationBeforeStart int                 `gorm:"default:15"` // in minutes
	AttendeeEmails          dbtypes.StringSlice `gorm:"type:jsonb"`

	Course *Course `gorm:"foreignKey:CourseID;references:ID"`
}

func (CourseEvent) TableName() string {
	return "course_events"
}
