package model

import "elearning-api/consts"

type User struct {
	BaseModel

	Email    string          `gorm:"type:varchar(255);not null;uniqueIndex"`
	Username string          `gorm:"type:varchar(100);not null;uniqueIndex"`
	Password string          `gorm:"type:text;not null"`
	Name     string          `gorm:"type:varchar(255)"`
	Avatar   string          `gorm:"type:text"`
	IsActive bool            `gorm:"default:true"`
	Role     consts.UserRole `gorm:"type:user_role;default:'student'"`
}
