package dto

type PagingRequest struct {
	Page  int `form:"page" json:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" json:"limit" binding:"omitempty,min=1,max=100"`
}

type UUIDRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (p *PagingRequest) Process() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 10
	}
}

func (p *PagingRequest) GetOffset() int {
	return (p.Page - 1) * p.Limit
}
