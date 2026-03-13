package dto

import (
	"elearning-api/model"
)

type InstructorProfileResponse struct {
	ID           string        `json:"id"`
	Bio          string        `json:"bio"`
	Education    string        `json:"education"`
	RatingAvg    float64       `json:"rating_avg"`
	TotalStudent int64         `json:"total_student"`
	TotalCourse  int64         `json:"total_course"`
	Balance      float64       `json:"balance"`
	LinkedinURL  string        `json:"linkedin_url"`
	YoutubeURL   string        `json:"youtube_url"`
	InstagramURL string        `json:"instagram_url"`
	User         *UserResponse `json:"user,omitempty"`

	Courses []*CourseResponse `json:"courses,omitempty"`
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
	if data.Courses != nil {
		courses = NewListCourseResponse(data.Courses)
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
		User:         user,
		Courses:      courses,
	}
}

type CreateInstructorProfileRequest struct {
	Bio          string  `json:"bio" binding:"required,min=3,max=1000"`
	Education    string  `json:"education" binding:"required,min=3,max=2000"`
	LinkedinURL  *string `json:"linkedin_url" binding:"omitempty,url"`
	YoutubeURL   *string `json:"youtube_url" binding:"omitempty,url"`
	InstagramURL *string `json:"instagram_url" binding:"omitempty,url"`
}

type UpdateInstructorProfileRequest struct {
	Bio          *string `json:"bio" binding:"omitempty,min=3,max=1000"`
	Education    *string `json:"education" binding:"omitempty,min=3,max=2000"`
	LinkedinURL  *string `json:"linkedin_url" binding:"omitempty,url"`
	YoutubeURL   *string `json:"youtube_url" binding:"omitempty,url"`
	InstagramURL *string `json:"instagram_url" binding:"omitempty,url"`
}

type ListInstructorProfileRequest struct {
	PagingRequest

	Name *string `form:"name"`
}
