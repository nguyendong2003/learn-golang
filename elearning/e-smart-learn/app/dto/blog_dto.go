package dto

import (
	"time"

	"elearning-api/model"

	"github.com/google/uuid"
)

type AuthorResponse struct {
	ID          string    `json:"id"`
	UserName    string    `json:"username"`
	Name        string    `json:"name"`
	Avatar      string    `json:"avatar"`
	CreatedAt   time.Time `json:"created_at"`
	IsFollowing bool      `json:"is_following"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BlogResponse struct {
	ID          uuid.UUID         `json:"id"`
	Slug        string            `json:"slug"`
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	ImageURL    string            `json:"image_url"`
	Author      *AuthorResponse   `json:"author"`
	ViewCount   int64             `json:"view_count"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Category    *CategoryResponse `json:"category"`
	Tags        []string          `json:"tags"`
	Status      string            `json:"status"`
	IsRead      bool              `json:"is_read"`
	ScheduledAt *time.Time        `json:"scheduled_at"`
	PublishedAt *time.Time        `json:"published_at"`
}

func NewListBlogResponse(blogs []*model.Blog) []*BlogResponse {
	res := make([]*BlogResponse, len(blogs))
	for i, b := range blogs {
		res[i] = NewBlogDetailResponse(b)
	}
	return res
}

func NewBlogDetailResponse(m *model.Blog) *BlogResponse {
	if m == nil {
		return nil
	}
	return &BlogResponse{
		ID:          m.ID,
		Slug:        m.Slug,
		Title:       m.Title,
		Content:     m.Content,
		ImageURL:    m.ImageURL,
		Author:      NewAuthorResponse(m.Author),
		ViewCount:   m.ViewTotal,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		Category:    NewCategoryDetailResponse(m.Category),
		Tags:        m.Tags,
		Status:      string(m.Status),
		IsRead:      false,
		ScheduledAt: m.ScheduledAt,
		PublishedAt: m.PublishedAt,
	}
}

func NewAuthorResponse(m *model.User) *AuthorResponse {
	if m == nil {
		return nil
	}
	return &AuthorResponse{
		ID:          m.ID.String(),
		UserName:    m.Username,
		Name:        m.Name,
		Avatar:      m.Avatar,
		IsFollowing: false,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

type CreateBlogRequest struct {
	Title       string     `json:"title" binding:"required"`
	CategoryID  uuid.UUID  `json:"category_id" binding:"required"`
	Content     string     `json:"content" binding:"required"`
	ImageURL    string     `json:"image_url" binding:"required"`
	Tags        []string   `json:"tags" binding:"omitempty,dive,required"`
	Status      string     `json:"status" binding:"required,oneof=draft published scheduled"`
	ScheduledAt *time.Time `json:"scheduled_at" binding:"omitempty,required_if=Status scheduled"`
}

type UpdateBlogRequest struct {
	Title       string     `json:"title" binding:"required"`
	CategoryID  uuid.UUID  `json:"category_id" binding:"required,uuid"`
	Content     string     `json:"content" binding:"required"`
	ImageURL    string     `json:"image_url" binding:"required"`
	Tags        []string   `json:"tags" binding:"omitempty,dive,required"`
	Status      string     `json:"status" binding:"required,oneof=draft published scheduled"`
	ScheduledAt *time.Time `json:"scheduled_at" binding:"omitempty,required_if=Status scheduled"`
}
type SearchBlogRequest struct {
	PagingRequest
	CategoryID string `form:"category_id" json:"category_id" binding:"omitempty,uuid"`
	Keyword    string `form:"keyword" json:"keyword" binding:"omitempty"`
	Status     string `form:"status" json:"status" binding:"omitempty,oneof=draft published scheduled"`
}

type BlogStatisticsResponse struct {
	TotalArticles int64   `json:"total_articles"`
	TotalViews    int64   `json:"total_views"`
	AvgEngagement float64 `json:"avg_engagement"`
	Published     int64   `json:"published"`
	Drafts        int64   `json:"drafts"`
	Scheduled     int64   `json:"scheduled"`
}
