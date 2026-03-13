package repository

import (
	"strings"

	"gorm.io/gorm"
)

type Preload string

const (
	User              Preload = "User"
	Author            Preload = "Author"
	InstructorProfile Preload = "InstructorProfile"
	Course            Preload = "Course"
	Courses           Preload = "Courses"
	Category          Preload = "Category"
	Role              Preload = "Role"
	Permissions       Preload = "Permissions"
)

// Use for nested preloads, e.g. "Author.User", "Course.InstructorProfile.User"
func PreloadPath(parts ...Preload) Preload {
	if len(parts) == 0 {
		return ""
	}

	strs := make([]string, len(parts))
	for i, p := range parts {
		strs[i] = string(p)
	}

	return Preload(strings.Join(strs, "."))
}

func applyPreloads(db *gorm.DB, preloads []Preload) *gorm.DB {
	for _, p := range preloads {
		db = db.Preload(string(p))
	}

	return db
}
