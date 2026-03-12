package util

import (
	"elearning-api/apperror"
	"elearning-api/consts"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
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

func GetRequestID(c *gin.Context) string {
	if v, ok := c.Get(consts.RequestIDKey); ok {
		if rid, ok := v.(string); ok {
			return rid
		}
	}
	return ""
}

func GenerateSlug(title string) string {
	return fmt.Sprintf("%s-%s", slug.Make(title), uuid.New().String()[:8])
}

func GetRequestUserID(c *gin.Context) (uuid.UUID, error) {
	var userID uuid.UUID
	if v, exists := c.Get(consts.ContextUserID); exists {
		if s, ok := v.(string); ok {
			if id, err := uuid.Parse(s); err == nil {
				userID = id
			} else {
				return uuid.Nil, apperror.NewUnauthorizedError("Invalid user ID in context")
			}
		}
	}
	return userID, nil
}
