package model

import (
	"elearning-api/consts"

	"github.com/google/uuid"
)

type Course struct {
	BaseModel
	Title        string              `gorm:"type:varchar(255);not null"`
	Description  string              `gorm:"type:text"`
	Image        string              `gorm:"type:text"`
	Slug         string              `gorm:"type:varchar(255);not null;uniqueIndex"`
	UserID       uuid.UUID           `gorm:"type:uuid;not null;index"`
	Duration     int                 `gorm:"default:0"`
	CategoryID   uuid.UUID           `gorm:"type:uuid;not null;index"`
	Price        float64             `gorm:"type:decimal(10,2);default:0"`
	OldPrice     float64             `gorm:"type:decimal(10,2);default:0"`
	AverageRate  float64             `gorm:"type:decimal(3,2);default:0"`
	Status       consts.CourseStatus `gorm:"type:course_status;default:'draft'"`
	TotalStudent int64               `gorm:"default:0"`

	User         *User          `gorm:"foreignKey:UserID;references:ID"`
	Category     *Category      `gorm:"foreignKey:CategoryID;references:ID"`
	Chapters     []*Chapter     `gorm:"foreignKey:CourseID;references:ID"`
	Enrollments  []*Enrollment  `gorm:"foreignKey:CourseID;references:ID"`
	Feedbacks    []*Feedback    `gorm:"foreignKey:CourseID;references:ID"`
	CourseEvents []*CourseEvent `gorm:"foreignKey:CourseID;references:ID"`
}

func (Course) TableName() string {
	return "courses"
}

type Chapter struct {
	BaseModel
	CourseID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`

	Course  *Course   `gorm:"foreignKey:CourseID;references:ID"`
	Lessons []*Lesson `gorm:"foreignKey:ChapterID;references:ID"`
}

func (Chapter) TableName() string {
	return "chapters"
}

type Lesson struct {
	BaseModel
	ChapterID       uuid.UUID `gorm:"type:uuid;not null;index"`
	Title           string    `gorm:"type:varchar(255);not null"`
	Content         string    `gorm:"type:text"`
	VideoURL        string    `gorm:"type:text"`
	IsAbleToPreview bool      `gorm:"default:false"`

	Chapter *Chapter `gorm:"foreignKey:ChapterID;references:ID"`
}

func (Lesson) TableName() string {
	return "lessons"
}
