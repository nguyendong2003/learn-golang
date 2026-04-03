package dto

import (
	"time"

	"elearning-api/consts"
	"elearning-api/model"

	"github.com/google/uuid"
)

type CourseResponse struct {
	ID                string                     `json:"id"`
	Title             string                     `json:"title"`
	Description       string                     `json:"description"`
	Image             string                     `json:"image"`
	Slug              string                     `json:"slug"`
	Duration          int                        `json:"duration"`
	Price             float64                    `json:"price"`
	OldPrice          float64                    `json:"old_price"`
	Status            consts.CourseStatus        `json:"status"`
	AverageRate       float64                    `json:"average_rate"`
	TotalStudent      int64                      `json:"total_student"`
	IsPurchased       bool                       `json:"is_purchased"`
	Category          *CategoryResponse          `json:"category,omitempty"`
	InstructorProfile *InstructorProfileResponse `json:"instructor,omitempty"`
	CreatedAt         time.Time                  `json:"created_at"`
	UpdatedAt         time.Time                  `json:"updated_at"`
}
type CourseStatisticsResponse struct {
	TotalCourses   int64 `json:"total_courses"`
	PendingReviews int64 `json:"pending_reviews"`
	Drafts         int64 `json:"drafts"`
	Published      int64 `json:"published"`
	Archived       int64 `json:"archived"`
}

type InstructorTaughtCourseRevenueResponse struct {
	CourseID     string              `json:"course_id"`
	Title        string              `json:"title"`
	Slug         string              `json:"slug"`
	Image        string              `json:"image"`
	Status       consts.CourseStatus `json:"status"`
	TotalStudent int64               `json:"total_student"`
	Revenue      float64             `json:"revenue"`
}

func NewInstructorTaughtCourseRevenueListResponse(courses []*model.InstructorTaughtCourseRevenue) []*InstructorTaughtCourseRevenueResponse {
	res := make([]*InstructorTaughtCourseRevenueResponse, len(courses))
	for i, c := range courses {
		if c == nil {
			continue
		}

		res[i] = &InstructorTaughtCourseRevenueResponse{
			CourseID:     c.CourseID.String(),
			Title:        c.Title,
			Slug:         c.Slug,
			Image:        c.Image,
			Status:       c.Status,
			TotalStudent: c.TotalStudent,
			Revenue:      float64(c.Revenue) / 100,
		}
	}

	return res
}

type NewCoursesLast30DaysResponse struct {
	TotalNewCourses int64 `json:"total_new_courses"`
}

func NewListCourseResponse(courses []*model.Course) []*CourseResponse {
	res := make([]*CourseResponse, len(courses))
	for i, c := range courses {
		res[i] = NewCourseDetailResponse(c)
	}
	return res
}

func NewCourseDetailResponse(m *model.Course) *CourseResponse {
	if m == nil {
		return nil
	}

	var category *CategoryResponse
	if m.Category != nil {
		category = NewCategoryDetailResponse(m.Category)
	}

	var instructorProfile *InstructorProfileResponse
	if m.User != nil && m.User.InstructorProfile != nil {
		instructorProfile = NewInstructorProfileDetailResponse(m.User.InstructorProfile)
	}

	return &CourseResponse{
		ID:                m.ID.String(),
		Title:             m.Title,
		Description:       m.Description,
		Image:             m.Image,
		Slug:              m.Slug,
		Duration:          m.Duration,
		Price:             m.Price,
		Status:            m.Status,
		AverageRate:       m.AverageRate,
		TotalStudent:      m.TotalStudent,
		IsPurchased:       false,
		Category:          category,
		InstructorProfile: instructorProfile,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

type CreateCourseRequest struct {
	Title       string    `json:"title" binding:"required,min=3,max=255"`
	Description string    `json:"description"`
	Price       float64   `json:"price" binding:"required,gt=0"`
	ImageURL    string    `json:"image_url" binding:"required,url"`
	CategoryID  uuid.UUID `json:"category_id" binding:"required,uuid"`
}

type UpdateCourseRequest struct {
	Title       *string  `json:"title" binding:"omitempty,min=3,max=255"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price" binding:"omitempty,gt=0"`
	ImageURL    *string  `json:"image_url" binding:"omitempty,url"`
	CategoryID  *string  `json:"category_id" binding:"omitempty,uuid"`
}

type UpdateCourseStatusRequest struct {
	Status consts.CourseStatus `json:"status" binding:"required,oneof=draft pending_review published rejected archived"`
}

type ListCourseRequest struct {
	PagingRequest

	Title      *string              `form:"title"`
	CategoryID *string              `form:"category_id" binding:"omitempty,uuid"`
	Status     *consts.CourseStatus `form:"status"`
	UserID     *string              `form:"-"`
}

type CourseSlugRequest struct {
	Slug string `uri:"slug" binding:"required"`
}

type CourseEventRequest struct {
	Name                    string    `json:"name" binding:"required,min=3,max=255"`
	StartTime               time.Time `json:"start_time" binding:"required"`
	EndTime                 time.Time `json:"end_time" binding:"required,gtfield=StartTime"`
	Location                string    `json:"location" binding:"required,min=3,max=255"`
	NotificationBeforeStart int       `json:"notification_before_start" binding:"required,gt=0"`
	AttendeeEmails          []string  `json:"attendee_emails" binding:"required,gt=0,dive,email"`
	Description             string    `json:"description"`
}
type CourseEventResponse struct {
	ID                 string    `json:"id"`
	CourseID           string    `json:"course_id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	Location           string    `json:"location"`
	StartTime          time.Time `json:"start_time"`
	EndTime            time.Time `json:"end_time"`
	NotificationBefore int       `json:"notification_before"`
	AttendeeEmails     []string  `json:"attendee_emails"`
	RoomToken          string    `json:"room_token"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func NewCourseEventResponse(m *model.CourseEvent) *CourseEventResponse {
	if m == nil {
		return nil
	}

	return &CourseEventResponse{
		ID:                 m.ID.String(),
		CourseID:           m.CourseID.String(),
		Name:               m.Name,
		Description:        m.Description,
		Location:           m.Location,
		StartTime:          m.StartTime,
		EndTime:            m.EndTime,
		NotificationBefore: m.NotificationBeforeStart,
		AttendeeEmails:     m.AttendeeEmails,
		RoomToken:          m.RoomToken,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

func NewListCourseEventResponse(events []*model.CourseEvent) []*CourseEventResponse {
	res := make([]*CourseEventResponse, len(events))
	for i, e := range events {
		res[i] = NewCourseEventResponse(e)
	}
	return res
}
