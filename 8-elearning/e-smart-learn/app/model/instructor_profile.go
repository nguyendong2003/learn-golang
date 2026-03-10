package model

import (
	"github.com/google/uuid"
)

type InstructorProfile struct {
	BaseModel
	UserID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	Bio          string    `gorm:"type:text"`
	Education    string    `gorm:"type:text"`
	RatingAvg    float64   `gorm:"type:decimal(3,2);default:0"`
	TotalStudent int64     `gorm:"default:0"`
	TotalCourse  int64     `gorm:"default:0"`
	Balance      float64   `gorm:"type:decimal(15,2);default:0"`
	LinkedinURL  string    `gorm:"type:varchar(255)"`
	YoutubeURL   string    `gorm:"type:varchar(255)"`
	InstagramURL string    `gorm:"type:varchar(255)"`

	User    *User     `gorm:"foreignKey:UserID;references:ID"`
	Courses []*Course `gorm:"foreignKey:InstructorID;references:ID"`
}

func (InstructorProfile) TableName() string {
	return "instructor_profiles"
}
