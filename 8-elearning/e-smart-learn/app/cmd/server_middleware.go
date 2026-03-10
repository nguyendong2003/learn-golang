package cmd

import (
	"bytes"
	"elearning-api/apperror"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/util"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *ApiServer) ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err

		if appErr, ok := err.(*apperror.AppError); ok {
			status := apperror.MapErrorToStatus(appErr.Type)

			res := dto.NewApiResponse(c)
			res.Status = dto.NewResponseStatus(status)
			res.Errors = []apperror.AppError{*appErr}
			res.Request = dto.GetRequestClient(c)

			c.JSON(status, res)

			return
		}

		// fallback
		c.JSON(http.StatusInternalServerError, apperror.NewInternalServerError("Something went wrong"))
	}
}

func (s *ApiServer) CaptureRequestBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body == nil {
			c.Next()
			return
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Next()
			return
		}

		// reset body để handler vẫn đọc được
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var body any
		if len(bodyBytes) > 0 {
			json.Unmarshal(bodyBytes, &body)
		}

		c.Set("request_body", body)

		c.Next()
	}
}

func (s *ApiServer) AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(apperror.NewUnauthorizedError("Missing authorization header"))
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		token, err := util.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.Error(apperror.NewUnauthorizedError(err.Error()))
			c.Abort()
			return
		}

		// Validate access token
		claims, err := util.ValidateAccessToken(token, &s.config.JWT)
		if err != nil {
			c.Error(apperror.NewUnauthorizedError("Invalid or expired token"))
			c.Abort()
			return
		}

		// Set user info in context for later use
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}

func (s *ApiServer) LoadRolePermissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Error(apperror.NewUnauthorizedError("User not authenticated"))
			c.Abort()
			return
		}

		var user model.User

		err := s.dbRepository.GetDB().
			Preload("Role").
			Preload("Role.Permissions").
			First(&user, "id = ?", userID).
			Error

		if err != nil {
			c.Error(apperror.NewUnauthorizedError("User not found"))
			c.Abort()
			return
		}

		// set role vào context
		c.Set("user_role", user.Role.Name)

		// set permissions vào context
		var permissions []string
		for _, p := range user.Role.Permissions {
			permissions = append(permissions, p.Code)
		}

		c.Set("user_permissions", permissions)

		c.Next()
	}
}

func (s *ApiServer) RequireRoleHandler(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("user_role")

		if userRole != role {
			c.Error(apperror.NewForbiddenError("Role " + role + " required"))
			c.Abort()
			return
		}

		c.Next()
	}
}

func (s *ApiServer) RequirePermissionHandler(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		perms, exists := c.Get("user_permissions")
		if !exists {
			c.Error(apperror.NewForbiddenError("Permissions not found"))
			c.Abort()
			return
		}

		permissions := perms.([]string)

		for _, p := range permissions {
			if p == permission {
				c.Next()
				return
			}
		}

		c.Error(apperror.NewForbiddenError("Permission denied"))
		c.Abort()
	}
}

func (s *ApiServer) CorsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func (s *ApiServer) RequestIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")

		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set(util.RequestIDKey, requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}
