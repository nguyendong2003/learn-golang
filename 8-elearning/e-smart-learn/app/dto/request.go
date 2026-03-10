package dto

type PagingRequest struct {
	Limit     int    `form:"limit" json:"limit" binding:"omitempty,min=1,max=100"`
	Offset    int    `form:"offset" json:"offset" binding:"omitempty,min=0"`
	SortBy    string `form:"sortBy" json:"sortBy" binding:"omitempty"`
	SortOrder string `form:"sortOrder" json:"sortOrder" binding:"omitempty,oneof=asc desc"`
}

type UUIDRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (p *PagingRequest) Process() {
	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 10
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
	if p.SortBy == "" {
		p.SortBy = "created_at"
	}
	if p.SortOrder == "" {
		p.SortOrder = "desc"
	}
}
