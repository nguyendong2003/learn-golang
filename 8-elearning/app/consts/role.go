package consts

type UserRole string

const (
	RoleAdmin      UserRole = "admin"
	RoleStudent    UserRole = "student"
	RoleInstructor UserRole = "instructor"
)
