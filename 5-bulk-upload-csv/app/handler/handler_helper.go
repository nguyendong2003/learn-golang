package handler

import (
	"bulk-upload-csv/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func errorResponse(c *gin.Context, response *dto.ApiResponse, errs []error) {
	response.Errors = errs

	statusCode := http.StatusBadRequest

	//

	//
	respStatus := dto.ApiResponseStatus{
		Code: statusCode,
		Type: http.StatusText(statusCode),
	}

	response.Status = respStatus

	c.JSON(statusCode, response)
}
