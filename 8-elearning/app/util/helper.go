package util

import (
	"elearning-api/apperror"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

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

func GetDuration(durationStr string) (time.Duration, error) {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format: %w", err)
	}
	return duration, nil
}

func MustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("%s is not set", key)
	}
	return val
}
