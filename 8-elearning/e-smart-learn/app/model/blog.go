package model

import "github.com/google/uuid"

type Blog struct {
	BaseModel
	Title     string    `gorm:"type:varchar(255);not null"`
	Content   string    `gorm:"type:text;not null"`
	Slug      string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	ImageURL     string    `gorm:"type:text;not null"`
	AuthorID  uuid.UUID `gorm:"type:uuid;not null;index"`
	ViewTotal int64     `gorm:"default:0"`

	Author *InstructorProfile `gorm:"foreignKey:AuthorID;references:UserID"`
}

func (Blog) TableName() string {
	return "blogs"
}
