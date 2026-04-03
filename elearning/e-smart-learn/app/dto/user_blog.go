package dto

type GetReadBlogsResponse struct {
	ReadBlogs []BlogResponse `json:"read_blogs"`
	Total     int            `json:"total"`
}

type ViewReadBlogHistoryRequest struct {
	PagingRequest
}

type UpdateReadBlogHistoryRequest struct {
	BlogID string `json:"blog_id" binding:"required,uuid"`
}
