package dto

type PagingRequest struct {
	// Dùng tag 'form' để Gin bind dữ liệu từ URL Query String
	// Dùng 'int' để tương thích tốt nhất với các hàm Limit/Offset của GORM
	Limit int `form:"limit,default=10" json:"limit" binding:"omitempty,min=1,max=100"`
	Page  int `form:"page,default=1" json:"page" binding:"omitempty,min=1"`
}

func (p *PagingRequest) GetOffset() int {
	return (p.Page - 1) * p.Limit
}
