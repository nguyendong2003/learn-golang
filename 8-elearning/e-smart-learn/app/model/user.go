package model

import "github.com/google/uuid"

type User struct {
	BaseModel

	Email    string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	Username string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	Password string    `gorm:"type:text;not null"`
	Name     string    `gorm:"type:varchar(255)"`
	Avatar   string    `gorm:"type:text"`
	IsActive bool      `gorm:"default:true"`
	RoleID   uuid.UUID `gorm:"type:uuid;not null;index"`

	Role              *Role              `gorm:"foreignKey:RoleID;references:ID"`
	RefreshTokens     []*RefreshToken    `gorm:"foreignKey:UserID;references:ID"`
	Blogs             []*Blog            `gorm:"foreignKey:AuthorID;references:ID"`
	Enrollments       []*Enrollment      `gorm:"foreignKey:UserID;references:ID"`
	Feedbacks         []*Feedback        `gorm:"foreignKey:UserID;references:ID"`
	InstructorProfile *InstructorProfile `gorm:"foreignKey:UserID;references:ID"`
	Subscriptions     []*Subscription    `gorm:"foreignKey:UserID;references:ID"`
	Courses           []*Course          `gorm:"foreignKey:UserID;references:ID"`
}
