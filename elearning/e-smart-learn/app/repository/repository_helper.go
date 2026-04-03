package repository

import (
	"strings"

	"gorm.io/gorm"
)

type Preload string

const (
	User              Preload = "User"
	Follower          Preload = "Follower"
	Followee          Preload = "Followee"
	Author            Preload = "Author"
	InstructorProfile Preload = "InstructorProfile"
	Course            Preload = "Course"
	Courses           Preload = "Courses"
	Category          Preload = "Category"
	Plan              Preload = "Plan"
	Subscriptions     Preload = "Subscriptions"
	Payments          Preload = "Payments"
	Role              Preload = "Role"
	Permissions       Preload = "Permissions"
	Chapters          Preload = "Chapters"
	Lessons           Preload = "Lessons"
	CourseEvents      Preload = "CourseEvents"
	Details           Preload = "Details" // For CoursePurchaseDetails field
)

type Join string

const (
	UserJoin Join = "User"
)

// Use for nested preloads, e.g. "Author.User", "Course.User.InstructorProfile"
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

// Use for nested joins, e.g. "Author", "Course.User"
func JoinPath(parts ...Join) Join {
	if len(parts) == 0 {
		return ""
	}

	strs := make([]string, len(parts))
	for i, p := range parts {
		strs[i] = string(p)
	}

	return Join(strings.Join(strs, "."))
}

func applyPreloads(db *gorm.DB, preloads []Preload) *gorm.DB {
	for _, p := range preloads {
		db = db.Preload(string(p))
	}

	return db
}

func applyJoins(db *gorm.DB, joins []Join) *gorm.DB {
	for _, j := range joins {
		db = db.Joins(string(j))
	}

	return db
}
