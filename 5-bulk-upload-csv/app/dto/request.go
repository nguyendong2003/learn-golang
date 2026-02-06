package dto

type PagingRequest struct {
	Limit  uint16 `query:"limit" json:"limit"`
	Offset uint16 `query:"offset" json:"offset"`
}
