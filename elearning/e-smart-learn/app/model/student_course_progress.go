package model

import (
	"time"

	"github.com/google/uuid"
)

type StudentCourseProgress struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID   uuid.UUID `gorm:"type:uuid;not null;index:idx_student_course_progress_user_course,priority:1"`
	CourseID uuid.UUID `gorm:"type:uuid;not null;index:idx_student_course_progress_user_course,priority:2"`
	LessonID uuid.UUID `gorm:"type:uuid;not null"`

	IsCompleted bool `gorm:"type:boolean;default:false"`
	CompletedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time

	User   *User   `gorm:"foreignKey:UserID;references:ID"`
	Course *Course `gorm:"foreignKey:CourseID;references:ID"`
	Lesson *Lesson `gorm:"foreignKey:LessonID;references:ID"`
}

func (StudentCourseProgress) TableName() string {
	return "student_course_progress"
}
