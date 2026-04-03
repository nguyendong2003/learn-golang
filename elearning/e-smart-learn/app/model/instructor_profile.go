package model

import (
	"elearning-api/consts"
	"elearning-api/repository/dbtypes"

	"github.com/google/uuid"
)

type InstructorProfile struct {
	BaseModel

	UserID            uuid.UUID                      `gorm:"type:uuid;not null;uniqueIndex"`
	CategoryID        uuid.UUID                      `gorm:"type:uuid;index"`
	Bio               string                         `gorm:"type:text"`
	Education         string                         `gorm:"type:text"`
	RatingAvg         float64                        `gorm:"type:decimal(3,2);default:0"`
	TotalStudent      int64                          `gorm:"default:0"`
	TotalCourse       int64                          `gorm:"default:0"`
	Balance           float64                        `gorm:"type:decimal(15,2);default:0"`
	LinkedinURL       string                         `gorm:"type:varchar(255)"`
	YoutubeURL        string                         `gorm:"type:varchar(255)"`
	InstagramURL      string                         `gorm:"type:varchar(255)"`
	YearsOfExperience int                            `gorm:"default:0"`
	CVURL             string                         `gorm:"type:text"`
	PortfolioURL      string                         `gorm:"type:text"`
	Certifications    dbtypes.StringSlice            `gorm:"type:jsonb"`
	Status            consts.InstructorProfileStatus `gorm:"type:text;default:'pending_review'"`

	User     *User     `gorm:"foreignKey:UserID;references:ID"`
	Category *Category `gorm:"foreignKey:CategoryID;references:ID"`
}

func (InstructorProfile) TableName() string {
	return "instructor_profiles"
}
