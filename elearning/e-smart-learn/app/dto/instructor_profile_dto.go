package dto

import (
	"elearning-api/model"
	"time"
)

type InstructorProfileResponse struct {
	ID           string            `json:"id"`
	Category     *CategoryResponse `json:"category,omitempty"`
	Bio          string            `json:"bio"`
	Education    string            `json:"education"`
	RatingAvg    float64           `json:"rating_avg"`
	TotalStudent int64             `json:"total_student"`
	TotalCourse  int64             `json:"total_course"`
	Balance      float64           `json:"balance"`
	LinkedinURL  string            `json:"linkedin_url"`
	YoutubeURL   string            `json:"youtube_url"`
	InstagramURL string            `json:"instagram_url"`
	User         *UserResponse     `json:"user,omitempty"`
	Courses      []*CourseResponse `json:"courses,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

func NewInstructorProfileListResponse(instructors []*model.InstructorProfile) []*InstructorProfileResponse {
	res := make([]*InstructorProfileResponse, len(instructors))
	for i, ins := range instructors {
		res[i] = NewInstructorProfileDetailResponse(ins)
	}
	return res
}

func NewInstructorProfileDetailResponse(data *model.InstructorProfile) *InstructorProfileResponse {
	if data == nil {
		return nil
	}

	var user *UserResponse
	if data.User != nil {
		user = NewUserDetailResponse(data.User)
	}

	var courses []*CourseResponse
	if data.User != nil && data.User.Courses != nil {
		courses = NewListCourseResponse(data.User.Courses)
	}

	var category *CategoryResponse
	if data.Category != nil {
		category = NewCategoryDetailResponse(data.Category)
	}

	return &InstructorProfileResponse{
		ID:           data.ID.String(),
		Bio:          data.Bio,
		Education:    data.Education,
		RatingAvg:    data.RatingAvg,
		TotalStudent: data.TotalStudent,
		TotalCourse:  data.TotalCourse,
		Balance:      data.Balance,
		LinkedinURL:  data.LinkedinURL,
		YoutubeURL:   data.YoutubeURL,
		InstagramURL: data.InstagramURL,
		Category:     category,
		User:         user,
		Courses:      courses,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
	}
}

type CreateInstructorProfileRequest struct {
	Bio          string  `json:"bio" binding:"required,min=3,max=1000"`
	Education    string  `json:"education" binding:"omitempty"`
	LinkedinURL  *string `json:"linkedin_url" binding:"omitempty,url"`
	YoutubeURL   *string `json:"youtube_url" binding:"omitempty,url"`
	InstagramURL *string `json:"instagram_url" binding:"omitempty,url"`
	CategoryID   *string `json:"category_id" binding:"omitempty,uuid"`
}

type UpdateInstructorProfileRequest struct {
	Bio          *string `json:"bio" binding:"omitempty,min=3,max=1000"`
	Education    *string `json:"education" binding:"omitempty,min=3,max=2000"`
	LinkedinURL  *string `json:"linkedin_url" binding:"omitempty,url"`
	YoutubeURL   *string `json:"youtube_url" binding:"omitempty,url"`
	InstagramURL *string `json:"instagram_url" binding:"omitempty,url"`
	CategoryID   *string `json:"category_id" binding:"omitempty,uuid"`
}

type ListInstructorProfileRequest struct {
	PagingRequest

	Name *string `form:"name"`
}

type ApplyInstructorRequest struct {
	Bio               string   `json:"bio" binding:"required,min=3,max=1000"`
	Education         string   `json:"education"`
	LinkedinURL       string   `json:"linkedin_url" binding:"omitempty,url"`
	YoutubeURL        string   `json:"youtube_url" binding:"omitempty,url"`
	InstagramURL      string   `json:"instagram_url" binding:"omitempty,url"`
	YearsOfExperience int      `json:"years_of_experience" binding:"required,gte=0"`
	CVURL             string   `json:"cv_url" binding:"omitempty,url"`
	PortfolioURL      string   `json:"portfolio_url" binding:"omitempty,url"`
	Certifications    []string `json:"certifications"`
	CategoryID        *string  `json:"category_id" binding:"omitempty,uuid"`
}

type InstructorApplicationResponse struct {
	ID string `json:"id"`
	ApplyInstructorRequest
	Status string        `json:"status"`
	User   *UserResponse `json:"user,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewInstructorApplicationListResponse(applications []*model.InstructorProfile) []*InstructorApplicationResponse {
	res := make([]*InstructorApplicationResponse, len(applications))
	for i, app := range applications {
		res[i] = NewInstructorApplicationResponse(app)
	}
	return res
}

func NewInstructorApplicationResponse(data *model.InstructorProfile) *InstructorApplicationResponse {
	if data == nil {
		return nil
	}

	var user *UserResponse
	if data.User != nil {
		user = NewUserDetailResponse(data.User)
	}

	return &InstructorApplicationResponse{
		ID: data.ID.String(),
		ApplyInstructorRequest: ApplyInstructorRequest{
			Bio:               data.Bio,
			Education:         data.Education,
			LinkedinURL:       data.LinkedinURL,
			YoutubeURL:        data.YoutubeURL,
			InstagramURL:      data.InstagramURL,
			YearsOfExperience: data.YearsOfExperience,
			CVURL:             data.CVURL,
			PortfolioURL:      data.PortfolioURL,
			Certifications:    data.Certifications,
		},
		Status:    string(data.Status),
		User:      user,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}
