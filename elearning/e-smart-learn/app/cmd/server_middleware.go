package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"elearning-api/apperror"
	"elearning-api/consts"
	"elearning-api/dto"
	"elearning-api/util"

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
			util.LoggerFromContextWithLayer(c.Request.Context(), util.LayerHTTP).Warn("request failed",
				"method", c.Request.Method,
				"path", c.FullPath(),
				"status", status,
				"error_type", appErr.Type,
				"error_message", appErr.Message,
			)

			res := dto.NewApiResponse(c)
			res.Status = dto.NewResponseStatus(status)
			res.Errors = []apperror.AppError{*appErr}
			res.Request = dto.GetRequestClient(c)

			c.JSON(status, res)

			return
		}

		// fallback
		util.LoggerFromContextWithLayer(c.Request.Context(), util.LayerHTTP).Error("unexpected request error",
			"method", c.Request.Method,
			"path", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusInternalServerError, apperror.NewInternalServerError("Something went wrong"))
	}
}

func (s *ApiServer) RequestLoggingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		status := c.Writer.Status()
		latencyMs := time.Since(start).Milliseconds()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		logger := util.LoggerFromContextWithLayer(c.Request.Context(), util.LayerHTTP)
		attrs := []any{
			"method", c.Request.Method,
			"path", path,
			"status", status,
			"latency_ms", latencyMs,
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"error_count", len(c.Errors),
		}

		switch {
		case status >= http.StatusInternalServerError:
			logger.Error("http request failed", attrs...)
		case status >= http.StatusBadRequest:
			logger.Warn("http request failed", attrs...)
		default:
			logger.Info("http request completed", attrs...)
		}
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
			if err := json.Unmarshal(bodyBytes, &body); err != nil {
				body = nil
			}
		}

		c.Set("request_body", body)

		c.Next()
	}
}

func (s *ApiServer) AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			_ = c.Error(apperror.NewUnauthorizedError("Missing authorization header"))
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		token, err := util.ExtractTokenFromHeader(authHeader)
		if err != nil {
			_ = c.Error(apperror.NewUnauthorizedError(err.Error()))
			c.Abort()
			return
		}

		// Validate access token
		claims, err := util.ValidateAccessToken(token, &s.config.JWT)
		if err != nil {
			_ = c.Error(apperror.NewUnauthorizedError("Invalid or expired token"))
			c.Abort()
			return
		}

		// Set user info in context for later use
		c.Set(consts.ContextUserID, claims.UserID)

		c.Next()
	}
}

// Allows requests WITHOUT authentication but will SET USER INFO in context if token is provided
func (s *ApiServer) OptionalAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		token, err := util.ExtractTokenFromHeader(authHeader)
		if err != nil {
			_ = c.Error(apperror.NewUnauthorizedError(err.Error()))
			c.Abort()
			return
		}

		claims, err := util.ValidateAccessToken(token, &s.config.JWT)
		if err != nil {
			_ = c.Error(apperror.NewUnauthorizedError("Invalid or expired token"))
			c.Abort()
			return
		}

		c.Set(consts.ContextUserID, claims.UserID)
		c.Next()
	}
}

func (s *ApiServer) LoadRolePermissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		user, err := s.userService.GetUserWithRole(c.Request.Context(), userID)
		if err != nil {
			_ = c.Error(apperror.NewUnauthorizedError("User not found"))
			c.Abort()
			return
		}
		c.Set(consts.ContextUserRole, user.Role.Name)

		var permissions []string
		for _, p := range user.Role.Permissions {
			permissions = append(permissions, p.Code)
		}

		c.Set(consts.ContextUserPermissions, permissions)

		c.Next()
	}
}

func (s *ApiServer) LoadRolePermissionOptionalHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get(consts.ContextUserID); !exists {
			c.Next()
			return
		}

		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		user, err := s.userService.GetUserWithRole(c.Request.Context(), userID)
		if err != nil {
			_ = c.Error(apperror.NewUnauthorizedError("User not found"))
			c.Abort()
			return
		}
		c.Set(consts.ContextUserRole, user.Role.Name)

		var permissions []string
		for _, p := range user.Role.Permissions {
			permissions = append(permissions, p.Code)
		}

		c.Set(consts.ContextUserPermissions, permissions)
		c.Next()
	}
}

func (s *ApiServer) RequireRoleHandler(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString(consts.ContextUserRole)

		if userRole != role {
			_ = c.Error(apperror.NewForbiddenError("Role " + role + " required"))
			c.Abort()
			return
		}

		c.Next()
	}
}

func (s *ApiServer) RequirePermissionHandler(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		perms, exists := c.Get(consts.ContextUserPermissions)
		if !exists {
			_ = c.Error(apperror.NewForbiddenError("Permissions not found"))
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

		_ = c.Error(apperror.NewForbiddenError("Permission denied"))
		c.Abort()
	}
}

func (s *ApiServer) RequireInstructorProfileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(apperror.NewUnauthorizedError("Invalid user ID"))
			c.Abort()
			return
		}

		instructor, err := s.instructorProfileService.GetByUserID(c.Request.Context(), userID)
		if err != nil {
			_ = c.Error(apperror.NewForbiddenError("Instructor profile required"))
			c.Abort()
			return
		}

		if instructor == nil {
			_ = c.Error(apperror.NewForbiddenError("Instructor profile required"))
			c.Abort()
			return
		}

		c.Set(consts.ContextInstructorProfile, instructor)

		c.Next()
	}
}

func (s *ApiServer) LoadCourseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var courseIDRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&courseIDRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			c.Abort()
			return
		}

		courseID, err := uuid.Parse(courseIDRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			c.Abort()
			return
		}

		// get course
		course, err := s.courseService.GetByID(c.Request.Context(), courseID)
		if err != nil {
			_ = c.Error(apperror.NewNotFoundError("Course not found"))
			c.Abort()
			return
		}

		c.Set(consts.ContextCourse, course)
		c.Next()
	}
}

func (s *ApiServer) RequireCourseOwnerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		// Get course from context
		courseRaw := util.GetCourse(c)
		course, ok := courseRaw.(*dto.CourseResponse)
		if !ok {
			_ = c.Error(apperror.NewInternalServerError("Can not get course"))
			c.Abort()
			return
		}

		if course == nil {
			_ = c.Error(apperror.NewNotFoundError("Course not found"))
			c.Abort()
			return
		}

		// check profile
		if course.InstructorProfile == nil || course.InstructorProfile.User == nil {
			_ = c.Error(apperror.NewInternalServerError("Instructor information is missing"))
			c.Abort()
			return
		}

		// check owner
		if course.InstructorProfile.User.ID != userID.String() {
			_ = c.Error(apperror.NewForbiddenError("You are not the owner of this course"))
			c.Abort()
			return
		}

		// c.Set(consts.ContextCourse, course)
		c.Next()
	}
}

func (s *ApiServer) RequireInstructorProfileOwnerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(apperror.NewUnauthorizedError("Invalid user ID"))
			c.Abort()
			return
		}

		// get instructor profile id from param
		instructorID := c.Param("id")
		if instructorID == "" {
			_ = c.Error(apperror.NewBadRequestError("Instructor profile ID is required"))
			c.Abort()
			return
		}

		// convert uuid
		iid, err := uuid.Parse(instructorID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid instructor profile ID"))
			c.Abort()
			return
		}

		// query instructor profile
		instructor, err := s.instructorProfileService.GetByID(c.Request.Context(), iid)
		if err != nil {
			_ = c.Error(apperror.NewNotFoundError("Instructor profile not found"))
			c.Abort()
			return
		}

		if instructor.User == nil {
			_ = c.Error(apperror.NewInternalServerError("User information is missing"))
			c.Abort()
			return
		}

		// check owner
		if instructor.User.ID != userID.String() {
			_ = c.Error(apperror.NewForbiddenError("You are not the owner of this instructor profile"))
			c.Abort()
			return
		}

		// lưu instructor vào context
		c.Set(consts.ContextInstructorProfile, instructor)

		c.Next()
	}
}

func (s *ApiServer) RequireActiveSubscriptionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		isActive, err := s.subscriptionService.HasActiveSubscription(c.Request.Context(), userID)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		if !isActive {
			_ = c.Error(apperror.NewForbiddenError("Active subscription required"))
			c.Abort()
			return
		}
		c.Next()
	}
}

func (s *ApiServer) CorsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, Authorization, X-Requested-With, X-CSRF-Token, ngrok-skip-browser-warning")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")

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

		c.Set(consts.RequestIDKey, requestID)
		ctx := util.WithLoggerAttrs(c.Request.Context(), "request_id", requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}
