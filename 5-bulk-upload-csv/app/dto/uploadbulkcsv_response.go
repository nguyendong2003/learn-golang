package dto

type RowError struct {
	Row     int    `json:"row"`
	Field   string `json:"field"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

type BulkUploadResponse struct {
	TotalProcessed int        `json:"total_processed"`
	TotalSuccess   int        `json:"total_success"`
	Errors         []RowError `json:"errors"`
}
