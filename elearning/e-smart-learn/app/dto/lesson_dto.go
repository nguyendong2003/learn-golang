package dto

import (
	"time"

	"elearning-api/model"
)

type LessonResponse struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Duration        int       `json:"duration"`
	VideoURL        string    `json:"video_url"`
	DocumentURL     string    `json:"document_url"`
	IsAbleToPreview bool      `json:"is_able_to_preview"`
	Order           int       `json:"order"`
	Type            string    `json:"type"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func NewListLessonResponse(lessons []*model.Lesson) []*LessonResponse {
	res := make([]*LessonResponse, len(lessons))
	for i, c := range lessons {
		res[i] = NewLessonDetailResponse(c)
	}
	return res
}

func NewLessonDetailResponse(data *model.Lesson) *LessonResponse {
	if data == nil {
		return nil
	}
	return &LessonResponse{
		ID:              data.ID.String(),
		Title:           data.Title,
		Duration:        data.Duration,
		VideoURL:        data.VideoURL,
		DocumentURL:     data.DocumentURL,
		IsAbleToPreview: data.IsAbleToPreview,
		Order:           data.Order,
		Type:            string(data.Type),
		CreatedAt:       data.CreatedAt,
		UpdatedAt:       data.UpdatedAt,
	}
}

type ChapterResponse struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	Order     int               `json:"order"`
	Lessons   []*LessonResponse `json:"lessons" binding:"omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

func NewChapterResponse(data *model.Chapter) *ChapterResponse {
	if data == nil {
		return nil
	}
	return &ChapterResponse{
		ID:        data.ID.String(),
		Title:     data.Title,
		Order:     data.Order,
		Lessons:   NewListLessonResponse(data.Lessons),
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func NewListChapterResponse(chapters []*model.Chapter) []*ChapterResponse {
	res := make([]*ChapterResponse, len(chapters))
	for i, c := range chapters {
		res[i] = NewChapterResponse(c)
	}
	return res
}

type CreateLessonsRequest struct {
	Chapters []CreateLessonRequest `json:"chapters" binding:"required,min=1,dive"`
}
type CreateLessonRequest struct {
	ChapterTitle string                     `json:"chapter_title" binding:"required,min=3,max=255"`
	Lessons      []CreateLessonAssetRequest `json:"lessons" binding:"required,min=1,dive"`
}
type CreateLessonAssetRequest struct {
	Title           string `json:"title" binding:"required,min=3,max=255"`
	Duration        int    `json:"duration" binding:"required,gt=0"`
	IsAbleToPreview bool   `json:"is_able_to_preview"`
	VideoURL        string `json:"video_url" binding:"omitempty,url"`
	DocumentURL     string `json:"document_url" binding:"omitempty,url"`
}
type UpdateCourseWithChaptersRequest struct {
	Chapters []UpdateChapterWithLessonsRequest `json:"chapters" binding:"required,min=1,dive"`
}
type UpdateChapterWithLessonsRequest struct {
	ID      string                         `json:"id" binding:"omitempty,uuid"` // Empty = new chapter, filled = update/keep existing
	Title   string                         `json:"title" binding:"required,min=3,max=255"`
	Lessons []UpdateLessonInChapterRequest `json:"lessons" binding:"required,min=1,dive"`
}
type UpdateLessonInChapterRequest struct {
	ID              string `json:"id" binding:"omitempty,uuid"` // Empty = new lesson, filled = update/keep existing
	Title           string `json:"title" binding:"required,min=3,max=255"`
	Duration        int    `json:"duration" binding:"required,gt=0"`
	IsAbleToPreview bool   `json:"is_able_to_preview"`
	VideoURL        string `json:"video_url" binding:"omitempty,url"`
	DocumentURL     string `json:"document_url" binding:"omitempty,url"`
}
