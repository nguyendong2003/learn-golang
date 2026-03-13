package consts

type CourseStatus string

const (
	CourseDraft     CourseStatus = "draft"
	CoursePending   CourseStatus = "pending_review"
	CoursePublished CourseStatus = "published"
	CourseRejected  CourseStatus = "rejected"
	CourseArchived  CourseStatus = "archived"
)
