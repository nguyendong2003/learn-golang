package util

import (
	"elearning-api/apperror"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindAndValidateJSON(c *gin.Context, obj any) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errors := make(map[string]string)
			for _, fe := range ve {
				errors[fe.Field()] = fe.Translate(GetTranslator())
			}
			return apperror.NewValidationError(errors)
		}
		return apperror.NewBadRequestError("Invalid JSON body")
	}

	return nil
}

// RequestID
const RequestIDKey = "request_id"

func GetRequestID(c *gin.Context) string {
	if v, ok := c.Get(RequestIDKey); ok {
		if rid, ok := v.(string); ok {
			return rid
		}
	}
	return ""
}
