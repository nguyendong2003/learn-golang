package util

import (
	"elearning-api/apperror"
	"elearning-api/consts"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

// Validate JSON request body
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

// Validate JSON string in form data
func BindAndValidateJSONFromString(jsonStr string, obj any) error {
	if err := json.Unmarshal([]byte(jsonStr), obj); err != nil {
		return apperror.NewBadRequestError("Invalid JSON data")
	}

	// get gin validator engine
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return apperror.NewInternalServerError("Validator engine not found")
	}

	if err := v.Struct(obj); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errs := make(map[string]string)
			for _, fe := range ve {
				errs[fe.Field()] = fe.Translate(GetTranslator())
			}
			return apperror.NewValidationError(errs)
		}
		return apperror.NewBadRequestError("Validation failed")
	}

	return nil
}

func BindAndValidateQuery(c *gin.Context, obj any) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errors := make(map[string]string)
			for _, fe := range ve {
				errors[fe.Field()] = fe.Translate(GetTranslator())
			}
			return apperror.NewValidationError(errors)
		}
		return apperror.NewBadRequestError("Invalid query parameters")
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

func GetCourse(c *gin.Context) any {
	if v, ok := c.Get(consts.ContextCourse); ok {
		return v
	}
	return nil
}

func GetChapter(c *gin.Context) any {
	if v, ok := c.Get(consts.ContextChapter); ok {
		return v
	}
	return nil
}

func GetFileTypeFromFilename(filename string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	imageExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	videoExts := map[string]bool{
		".mp4":  true,
		".mov":  true,
		".avi":  true,
		".mkv":  true,
		".webm": true,
	}

	if imageExts[ext] {
		return "images", nil
	}

	if videoExts[ext] {
		return "videos", nil
	}

	return "", apperror.NewBadRequestError("Unsupported file type")
}

func ExtractFileURLFromPresign(presignURL string) (string, error) {
	parsed, err := url.Parse(presignURL)
	if err != nil {
		return "", apperror.NewBadRequestError("Invalid presigned URL")
	}

	// delete query params
	parsed.RawQuery = ""

	return parsed.String(), nil
}

func GetRole(c *gin.Context) string {
	if v, exists := c.Get(consts.ContextUserRole); exists {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
