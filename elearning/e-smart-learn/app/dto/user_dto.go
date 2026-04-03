package dto

import (
	"time"

	"elearning-api/model"
	"elearning-api/repository/dbtypes"
)

type UserResponse struct {
	ID                    string             `json:"id"`
	Email                 string             `json:"email"`
	Username              string             `json:"username"`
	Name                  string             `json:"name"`
	PhoneNumber           string             `json:"phone_number"`
	Address               string             `json:"address"`
	Avatar                string             `json:"avatar"`
	CreatedAt             time.Time          `json:"created_at"`
	UpdatedAt             time.Time          `json:"updated_at"`
	Role                  *RoleResponse      `json:"role,omitempty"`
	ReadBlogIDs           []string           `json:"read_blog_ids,omitempty"`
	InstructorProfile     *InstructorProfile `json:"instructor_profile,omitempty"`
	TotalCoursesEnrolled  int64              `json:"total_courses_enrolled"`
	TotalCourseInProgress int64              `json:"total_course_in_progress"`
}
type InstructorProfile struct {
	ID                string            `json:"id"`
	Bio               string            `json:"bio"`
	Education         string            `json:"education"`
	RatingAvg         float64           `json:"rating_avg"`
	TotalStudent      int64             `json:"total_student"`
	TotalCourse       int64             `json:"total_course"`
	Balance           float64           `json:"balance"`
	LinkedinURL       string            `json:"linkedin_url"`
	YoutubeURL        string            `json:"youtube_url"`
	InstagramURL      string            `json:"instagram_url"`
	YearsOfExperience int               `json:"years_of_experience"`
	CVURL             string            `json:"cv_url"`
	PortfolioURL      string            `json:"portfolio_url"`
	Certifications    []string          `json:"certifications"`
	Status            string            `json:"status"`
	Category          *CategoryResponse `json:"category,omitempty"`
}

func NewUserDetailResponse(data *model.User) *UserResponse {
	if data == nil {
		return nil
	}
	return &UserResponse{
		ID:                data.ID.String(),
		Email:             data.Email,
		Username:          data.Username,
		Name:              data.Name,
		Avatar:            data.Avatar,
		PhoneNumber:       data.PhoneNumber,
		Address:           data.Address,
		CreatedAt:         data.CreatedAt,
		UpdatedAt:         data.UpdatedAt,
		Role:              NewRoleResponse(data.Role),
		ReadBlogIDs:       dbtypes.ToStringSlice(data.ReadBlogIDs),
		InstructorProfile: NewInstructorProfileResponse(data.InstructorProfile),
	}
}

func NewInstructorProfileResponse(data *model.InstructorProfile) *InstructorProfile {
	if data == nil {
		return nil
	}
	return &InstructorProfile{
		ID:                data.ID.String(),
		Bio:               data.Bio,
		Education:         data.Education,
		RatingAvg:         data.RatingAvg,
		TotalStudent:      data.TotalStudent,
		TotalCourse:       data.TotalCourse,
		Balance:           data.Balance,
		LinkedinURL:       data.LinkedinURL,
		YoutubeURL:        data.YoutubeURL,
		InstagramURL:      data.InstagramURL,
		YearsOfExperience: data.YearsOfExperience,
		CVURL:             data.CVURL,
		PortfolioURL:      data.PortfolioURL,
		Certifications:    data.Certifications,
		Status:            string(data.Status),
		Category:          NewCategoryDetailResponse(data.Category),
	}
}

func NewListUserResponse(users []*model.User) []*UserResponse {
	res := make([]*UserResponse, len(users))
	for i, u := range users {
		res[i] = NewUserDetailResponse(u)
	}
	return res
}

type UserStatisticsResponse struct {
	TotalActiveUsers        int64   `json:"total_active_users"`
	ActiveUsersGrowthPct    float64 `json:"active_users_growth_pct"`
	PendingInstructors      int64   `json:"pending_instructors"`
	CourseCompletionRatePct float64 `json:"course_completion_rate_pct"`
}

type ActiveStudentsStatisticsResponse struct {
	TotalActiveStudents     int64   `json:"total_active_students"`
	ActiveStudentsGrowthPct float64 `json:"active_students_growth_pct"`
}

type UserListResponse struct {
	ID                              string  `json:"id"`
	Name                            string  `json:"name"`
	Email                           string  `json:"email"`
	Avatar                          string  `json:"avatar"`
	Role                            string  `json:"role"`
	Status                          string  `json:"status"`
	HasPendingInstructorApplication bool    `json:"has_pending_instructor_application"`
	ActiveCourses                   int64   `json:"active_courses"`
	TotalCoursesTaught              int64   `json:"total_courses_taught"`
	CompletedLessons                int64   `json:"completed_lessons"`
	TotalLessons                    int64   `json:"total_lessons"`
	OverallProgressPct              float64 `json:"overall_progress_pct"`
}

type UpdateUserRequest struct {
	Name        string `json:"name" binding:"omitempty"`
	Email       string `json:"email" binding:"omitempty,email"`
	PhoneNumber string `json:"phone_number" binding:"omitempty,e164"`
	Address     string `json:"address" binding:"omitempty"`
	Avatar      string `json:"avatar" binding:"omitempty,url"`
}
