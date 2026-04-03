package model

import (
	"elearning-api/repository/dbtypes"

	"github.com/google/uuid"
)

type User struct {
	BaseModel

	Email            string            `gorm:"type:varchar(255);not null;uniqueIndex"`
	Username         string            `gorm:"type:varchar(100);not null;uniqueIndex"`
	Password         string            `gorm:"type:text;not null"`
	Name             string            `gorm:"type:varchar(255)"`
	Avatar           string            `gorm:"type:text"`
	OauthProvider    *string           `gorm:"type:varchar(50);index"`
	OauthID          *string           `gorm:"type:varchar(255);index"`
	PhoneNumber      string            `gorm:"type:varchar(20)"`
	Address          string            `gorm:"type:text"`
	IsActive         bool              `gorm:"default:true"`
	StripeCustomerID *string           `gorm:"type:varchar(255);uniqueIndex"`
	RoleID           uuid.UUID         `gorm:"type:uuid;not null;index"`
	ReadBlogIDs      dbtypes.UUIDSlice `gorm:"type:uuid[]"`

	Role              *Role              `gorm:"foreignKey:RoleID;references:ID"`
	RefreshTokens     []*RefreshToken    `gorm:"foreignKey:UserID;references:ID"`
	Blogs             []*Blog            `gorm:"foreignKey:AuthorID;references:ID"`
	Enrollments       []*Enrollment      `gorm:"foreignKey:UserID;references:ID"`
	Feedbacks         []*Feedback        `gorm:"foreignKey:UserID;references:ID"`
	InstructorProfile *InstructorProfile `gorm:"foreignKey:UserID;references:ID"`
	Subscriptions     []*Subscription    `gorm:"foreignKey:UserID;references:ID"`
	Courses           []*Course          `gorm:"foreignKey:UserID;references:ID"`
}

type UserDirectoryRow struct {
	UserID                          string `gorm:"column:user_id"`
	Name                            string `gorm:"column:name"`
	Email                           string `gorm:"column:email"`
	Avatar                          string `gorm:"column:avatar"`
	RoleName                        string `gorm:"column:role_name"`
	IsActive                        bool   `gorm:"column:is_active"`
	HasPendingInstructorApplication bool   `gorm:"column:has_pending_instructor_application"`
	ActiveCourses                   int64  `gorm:"column:active_courses"`
	TotalCoursesTaught              int64  `gorm:"column:total_courses_taught"`
	CompletedLessons                int64  `gorm:"column:completed_lessons"`
	TotalLessons                    int64  `gorm:"column:total_lessons"`
}
