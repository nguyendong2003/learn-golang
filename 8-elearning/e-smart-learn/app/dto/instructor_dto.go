package dto

import "elearning-api/model"

type InstructorResponse struct {
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
}

func NewInstructorListResponse(instructors []*model.InstructorProfile) []*InstructorResponse {
	res := make([]*InstructorResponse, len(instructors))
	for i, ins := range instructors {
		res[i] = NewInstructorResponse(ins)
	}
	return res
}

func NewInstructorResponse(data *model.InstructorProfile) *InstructorResponse {
	var user *UserResponse
	if data.User != nil {
		user = NewUserDetailResponse(data.User)
	}
	
	return &InstructorResponse{
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
	}
}