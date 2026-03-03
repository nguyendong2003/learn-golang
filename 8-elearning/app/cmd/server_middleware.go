package cmd

import (
	"elearning-api/apperror"
	"elearning-api/consts"
	"elearning-api/dto"
	"elearning-api/util"
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
		requestID := util.GetRequestID(c)

		status := http.StatusInternalServerError
		code := apperror.InternalServer
		message := "Internal server error"

		if appErr, ok := err.(*apperror.AppError); ok {
			code = appErr.Type
			message = appErr.Message

			switch appErr.Type {
			case apperror.NotFound:
				status = http.StatusNotFound
			case apperror.BadRequest:
				status = http.StatusBadRequest
			case apperror.Unauthorized:
				status = http.StatusUnauthorized
			case apperror.Forbidden:
				status = http.StatusForbidden
			case apperror.DuplicateEntry:
				status = http.StatusConflict
			}
		}

		resp := dto.NewErrorResponse(
			string(code),
			message,
			requestID,
		)

		c.AbortWithStatusJSON(status, resp)
	}
}

func (s *ApiServer) AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
		// Implement token validation logic here
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
func (s *ApiServer) AuthorizationHandler(role consts.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Example authorization logic
		userRole := c.GetString("user_role")
		if userRole != string(role) {
			c.Error(apperror.NewForbiddenError("User does not have the required role"))
			return
		}
		c.Next()
	}
}

// Middleware để gán request ID cho mỗi request, giúp dễ dàng theo dõi log và debug
func (s *ApiServer) RequestIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Nếu client đã gửi X-Request-ID thì dùng lại
		requestID := c.GetHeader("X-Request-ID")

		if requestID == "" {
			requestID = uuid.NewString()
		}

		// Lưu vào context để handler khác dùng
		c.Set(util.RequestIDKey, requestID)

		// Trả lại cho client qua header
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}
