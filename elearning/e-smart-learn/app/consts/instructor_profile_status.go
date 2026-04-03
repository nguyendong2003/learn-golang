package consts

type InstructorProfileStatus string

const (
	InstructorProfilePending  InstructorProfileStatus = "pending_review"
	InstructorProfileApproved InstructorProfileStatus = "approved"
	InstructorProfileRejected InstructorProfileStatus = "rejected"
)
