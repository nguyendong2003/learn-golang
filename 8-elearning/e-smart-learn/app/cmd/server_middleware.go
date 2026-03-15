package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

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
		c.Set(consts.ContextUserID, claims.UserID)

		c.Next()
	}
}

func (s *ApiServer) LoadRolePermissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			c.Error(apperror.NewUnauthorizedError("Invalid user ID"))
			c.Abort()
			return
		}

		user, err := s.userService.GetUserWithRole(c.Request.Context(), userID)
		if err != nil {
			c.Error(apperror.NewUnauthorizedError("User not found"))
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

func (s *ApiServer) RequireInstructorProfileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, err := util.GetRequestUserID(c)
		if err != nil {
			c.Error(apperror.NewUnauthorizedError("Invalid user ID"))
			c.Abort()
			return
		}

		instructor, err := s.instructorService.GetByUserID(c.Request.Context(), userID)
		if err != nil {
			c.Error(apperror.NewForbiddenError("Instructor profile required"))
			c.Abort()
			return
		}

		if instructor == nil {
			c.Error(apperror.NewForbiddenError("Instructor profile required"))
			c.Abort()
			return
		}

		// lưu instructor vào context để handler dùng
		c.Set(consts.ContextInstructorProfile, instructor)

		c.Next()
	}
}

func (s *ApiServer) RequireCourseOwnerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		// get course id from param
		courseID := c.Param("id")
		if courseID == "" {
			c.Error(apperror.NewBadRequestError("Course ID is required"))
			c.Abort()
			return
		}

		// convert uuid
		cid, err := uuid.Parse(courseID)
		if err != nil {
			c.Error(apperror.NewBadRequestError("Invalid course ID"))
			c.Abort()
			return
		}

		// get course
		course, err := s.courseService.GetByID(c.Request.Context(), cid)
		if err != nil {
			c.Error(apperror.NewNotFoundError("Course not found"))
			c.Abort()
			return
		}

		if course.Instructor == nil || course.Instructor.User == nil {
			c.Error(apperror.NewInternalServerError("Instructor information is missing"))
			c.Abort()
			return
		}

		// check owner
		if course.Instructor.User.ID != userID.String() {
			c.Error(apperror.NewForbiddenError("You are not the owner of this course"))
			c.Abort()
			return
		}

		// save course to context for later use
		c.Set(consts.ContextCourse, course)

		c.Next()
	}
}

func (s *ApiServer) RequireInstructorProfileOwnerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			c.Error(apperror.NewUnauthorizedError("Invalid user ID"))
			c.Abort()
			return
		}

		// get instructor profile id from param
		instructorID := c.Param("id")
		if instructorID == "" {
			c.Error(apperror.NewBadRequestError("Instructor profile ID is required"))
			c.Abort()
			return
		}

		// convert uuid
		iid, err := uuid.Parse(instructorID)
		if err != nil {
			c.Error(apperror.NewBadRequestError("Invalid instructor profile ID"))
			c.Abort()
			return
		}

		// query instructor profile
		instructor, err := s.instructorService.GetByID(c.Request.Context(), iid)
		if err != nil {
			c.Error(apperror.NewNotFoundError("Instructor profile not found"))
			c.Abort()
			return
		}

		if instructor.User == nil {
			c.Error(apperror.NewInternalServerError("User information is missing"))
			c.Abort()
			return
		}

		// check owner
		if instructor.User.ID != userID.String() {
			c.Error(apperror.NewForbiddenError("You are not the owner of this instructor profile"))
			c.Abort()
			return
		}

		// lưu instructor vào context
		c.Set(consts.ContextInstructorProfile, instructor)

		c.Next()
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

		c.Set(consts.RequestIDKey, requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}
