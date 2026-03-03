package dto

// Dùng tag 'form' để Gin bind dữ liệu từ URL Query String
// Dùng 'int' để tương thích tốt nhất với các hàm Limit/Offset của GORM
type PagingRequest struct {
	Page     int `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"`
}

func (p *PagingRequest) Process() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 || p.PageSize > 100 {
		p.PageSize = 10
	}
}

func (p *PagingRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *PagingRequest) GetLimit() int {
	return p.PageSize
}
