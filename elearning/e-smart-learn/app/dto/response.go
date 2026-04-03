package dto

import (
	"net/http"

	"elearning-api/apperror"
	"elearning-api/util"

	"github.com/gin-gonic/gin"
)

type ResponseStatus struct {
	Code int    `json:"code"`
	Type string `json:"type"`
}

type Pagination struct {
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"`
	Total     int    `json:"total"`
}

type ApiResponse struct {
	ProcessID string              `json:"process_id"`
	Path      string              `json:"path"`
	Status    ResponseStatus      `json:"status"`
	Request   any                 `json:"request"`
	Errors    []apperror.AppError `json:"errors"`
	Data      any                 `json:"data"`
	Metadata  any                 `json:"metadata"`
}

type RequestClient struct {
	Params map[string]string `json:"params"`
	Query  map[string]any    `json:"query"`
	Body   map[string]any    `json:"body"`
}

type PresignUploadURLResponse struct {
	URL string `json:"url"`
}

func GetRequestClient(c *gin.Context) RequestClient {
	// path params
	params := make(map[string]string)
	for _, p := range c.Params {
		params[p.Key] = p.Value
	}

	// query params
	query := make(map[string]any)
	for k, v := range c.Request.URL.Query() {
		if len(v) == 1 {
			query[k] = v[0]
		} else {
			query[k] = v
		}
	}

	return RequestClient{
		Params: params,
		Query:  query,
		Body:   GetRequestBody(c),
	}
}

func GetRequestBody(c *gin.Context) map[string]any {
	if body, exists := c.Get("request_body"); exists {
		if bodyMap, ok := body.(map[string]any); ok {
			return bodyMap
		}
	}
	return map[string]any{}
}

func NewResponseStatus(code int) ResponseStatus {
	return ResponseStatus{
		Code: code,
		Type: http.StatusText(code),
	}
}

func NewPagination(limit, offset, total int, sortBy, sortOrder string) Pagination {
	return Pagination{
		Limit:     limit,
		Offset:    offset,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Total:     total,
	}
}

func NewApiResponse(c *gin.Context) *ApiResponse {
	processID := util.GetRequestID(c)
	path := c.Request.URL.Path

	respStatus := NewResponseStatus(http.StatusOK)

	respError := []apperror.AppError{}
	respData := make([]any, 0)
	respMetadata := make(map[string]any)

	return &ApiResponse{
		ProcessID: processID,
		Status:    respStatus,
		Errors:    respError,
		Data:      respData,
		Path:      path,
		Metadata:  respMetadata,
	}
}
