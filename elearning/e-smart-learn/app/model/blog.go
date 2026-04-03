package model

import (
	"time"

	"elearning-api/consts"
	"elearning-api/repository/dbtypes"

	"github.com/google/uuid"
)

type Blog struct {
	BaseModel
	Title       string              `gorm:"type:varchar(255);not null"`
	Content     string              `gorm:"type:text;not null"`
	Slug        string              `gorm:"type:varchar(255);not null;uniqueIndex"`
	CategoryID  uuid.UUID           `gorm:"type:uuid;not null;index"`
	ImageURL    string              `gorm:"type:text;not null"`
	AuthorID    uuid.UUID           `gorm:"type:uuid;not null;index"`
	ViewTotal   int64               `gorm:"default:0"`
	Tags        dbtypes.StringSlice `gorm:"type:jsonb"`
	Status      consts.BlogStatus   `gorm:"type:varchar(20);not null;default:'draft'"`
	ScheduledAt *time.Time          `gorm:"type:timestamp with time zone"`
	PublishedAt *time.Time          `gorm:"type:timestamp with time zone"`

	Author   *User     `gorm:"foreignKey:AuthorID;references:ID"`
	Category *Category `gorm:"foreignKey:CategoryID;references:ID"`
}

func (Blog) TableName() string {
	return "blogs"
}

type BlogStats struct {
	TotalArticles int64
	TotalViews    int64
	Published     int64
	Drafts        int64
	Scheduled     int64
	Engaged       int64
}
