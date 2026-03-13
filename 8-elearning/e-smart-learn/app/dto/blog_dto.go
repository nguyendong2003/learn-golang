package dto

import (
	"elearning-api/model"
)

type BlogResponse struct {
	ID        string                     `json:"id"`
	Title     string                     `json:"title"`
	Content   string                     `json:"content"`
	ImageURL  string                     `json:"image_url"`
	Author    *InstructorProfileResponse `json:"author"`
	ViewCount int64                      `json:"view_count"`
	CreatedAt string                     `json:"created_at"`
	UpdatedAt string                     `json:"updated_at"`
}

func NewListBlogResponse(blogs []*model.Blog) []*BlogResponse {
	res := make([]*BlogResponse, len(blogs))
	for i, b := range blogs {
		res[i] = NewBlogDetailResponse(b)
	}
	return res
}

func NewBlogDetailResponse(m *model.Blog) *BlogResponse {
	var author *InstructorProfileResponse
	if m.Author != nil {
		author = NewInstructorProfileDetailResponse(m.Author)
	}

	return &BlogResponse{
		ID:        m.Slug,
		Title:     m.Title,
		Content:   m.Content,
		ImageURL:  m.ImageURL,
		Author:    author,
		ViewCount: m.ViewTotal,
		CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: m.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

type CreateBlogRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	ImageURL string `json:"image_url" binding:"required"`
}
