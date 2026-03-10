package consts

type UserRole string

const (
	RoleAdmin      UserRole = "admin"
	RoleInstructor UserRole = "instructor"
	RoleStudent    UserRole = "student"
)
