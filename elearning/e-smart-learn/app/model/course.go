package model

import (
	"time"

	"elearning-api/consts"

	"github.com/google/uuid"
)

type Course struct {
	BaseModel
	Title           string              `gorm:"type:varchar(255);not null"`
	Description     string              `gorm:"type:text"`
	Image           string              `gorm:"type:text"`
	Slug            string              `gorm:"type:varchar(255);not null;uniqueIndex"`
	UserID          uuid.UUID           `gorm:"type:uuid;not null;index"`
	Duration        int                 `gorm:"default:0"`
	CategoryID      uuid.UUID           `gorm:"type:uuid;not null;index"`
	Price           float64             `gorm:"type:decimal(10,2);default:0"`
	AverageRate     float64             `gorm:"type:decimal(3,2);default:0"`
	Status          consts.CourseStatus `gorm:"type:course_status;default:'draft'"`
	TotalStudent    int64               `gorm:"default:0"`
	StripeProductID string              `gorm:"type:varchar(255);index"`
	StripePriceID   string              `gorm:"type:varchar(255);index"`
	StripeCurrency  string              `gorm:"type:varchar(10);default:'usd'"`
	StripeAmount    int64               `gorm:"default:0"`
	StripeSyncedAt  *time.Time

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

type CourseStatistics struct {
	TotalCourses   int64
	PendingReviews int64
	Drafts         int64
	Published      int64
	Archived       int64
}

type InstructorTaughtCourseRevenue struct {
	CourseID     uuid.UUID
	Title        string
	Slug         string
	Image        string
	Status       consts.CourseStatus
	TotalStudent int64
	Revenue      int64
	CreatedAt    time.Time
}

type Chapter struct {
	BaseModel
	CourseID uuid.UUID `gorm:"type:uuid;not null;index"`
	Title    string    `gorm:"type:varchar(255);not null"`
	Order    int       `gorm:"default:0"`
	Course   *Course   `gorm:"foreignKey:CourseID;references:ID"`
	Lessons  []*Lesson `gorm:"foreignKey:ChapterID;references:ID"`
}

func (Chapter) TableName() string {
	return "chapters"
}

type Lesson struct {
	BaseModel
	ChapterID       uuid.UUID         `gorm:"type:uuid;not null;index"`
	Title           string            `gorm:"type:varchar(255);not null"`
	Duration        int               `gorm:"default:0"`
	VideoURL        string            `gorm:"type:text"`
	DocumentURL     string            `gorm:"type:text"`
	IsAbleToPreview bool              `gorm:"default:false"`
	Order           int               `gorm:"default:0"`
	Type            consts.LessonType `gorm:"type:lesson_type;default:'video'"`
	Chapter         *Chapter          `gorm:"foreignKey:ChapterID;references:ID"`
}

func (Lesson) TableName() string {
	return "lessons"
}
