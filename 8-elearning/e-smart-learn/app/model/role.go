package model

import (
	"github.com/google/uuid"
)

type Role struct {
	BaseModel
	Name        string `gorm:"type:varchar(100);not null;uniqueIndex"`
	Description string `gorm:"type:text"`

	Users       []*User       `gorm:"foreignKey:RoleID;references:ID"`
	Permissions []*Permission `gorm:"many2many:role_permissions;"`
}

func (Role) TableName() string {
	return "roles"
}

type Permission struct {
	BaseModel
	Code        string `gorm:"type:varchar(100);not null;uniqueIndex"`
	Description string `gorm:"type:text"`

	Roles []*Role `gorm:"many2many:role_permissions;"`
}

func (Permission) TableName() string {
	return "permissions"
}

type RolePermission struct {
	RoleID       uuid.UUID `gorm:"type:uuid;not null;primaryKey;index"`
	PermissionID uuid.UUID `gorm:"type:uuid;not null;primaryKey;index"`

	Role       *Role       `gorm:"foreignKey:RoleID;references:ID"`
	Permission *Permission `gorm:"foreignKey:PermissionID;references:ID"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
