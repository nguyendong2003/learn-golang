package model

import (
	"time"

	"elearning-api/consts"
)

type PresignedUploadTracking struct {
	BaseModel
	ObjectURL  string                       `gorm:"type:text;unique;not null"`
	ObjectName string                       `gorm:"type:text;not null"`
	Filetype   string                       `gorm:"type:varchar(20);not null"`
	Status     consts.PresignedUploadStatus `gorm:"type:varchar(20);not null;default:'pending'"`
	ExpiresAt  time.Time
}
