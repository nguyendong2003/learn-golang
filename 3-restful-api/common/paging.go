package common

// Phân trang (dùng tag form để parse query params từ client truyền lên)
type Paging struct {
	Page  int   `json:"page" form:"page"`
	Limit int   `json:"limit" form:"limit"`
	Total int64 `json:"total" form:"-"` // Không nhận trường total từ client truyền lên
}

// Xử lí query params client truyền lên
func (p *Paging) Process() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 || p.Limit >= 100 {
		p.Limit = 10
	}
}
