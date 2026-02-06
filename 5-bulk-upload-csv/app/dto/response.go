package dto

import (
	"encoding/json"
	"net/http"
)

type ApiResponseStatus struct {
	Code int    `json:"code"`
	Type string `json:"type"`
}

type ApiResponse struct {
	Path     string            `json:"path"`
	Status   ApiResponseStatus `json:"status"`
	Request  any               `json:"request"`
	Errors   []error           `json:"errors"`
	Data     any               `json:"data"`
	Metadata any               `json:"metadata"`
}

func (response *ApiResponse) ToJSON() string {
	jsonout, err := json.Marshal(response)
	if err != nil {
		return ""
	}

	return string(jsonout)
}

func NewApiResponse(path string) *ApiResponse {
	respStatus := ApiResponseStatus{
		Code: http.StatusOK,
		Type: http.StatusText(http.StatusOK),
	}

	respError := []error{}
	respData := make([]any, 0)

	return &ApiResponse{
		Status: respStatus,
		Errors: respError,
		Data:   respData,
		Path:   path,
	}
}
