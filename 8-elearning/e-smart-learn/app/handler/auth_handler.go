package handler

import (
	"net/http"

	"elearning-api/apperror"
	"elearning-api/consts"
	"elearning-api/dto"
	"elearning-api/service"
	"elearning-api/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler interface {
	Login() gin.HandlerFunc
	Register() gin.HandlerFunc
	RefreshToken() gin.HandlerFunc
	ChangePassword() gin.HandlerFunc
	ForgotPassword() gin.HandlerFunc
	ResetPassword() gin.HandlerFunc
}

type authHandler struct {
	authService service.AuthService
}

func NewAuthHandler(
	authService service.AuthService,
) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

// Login godoc
// @Summary Login
// @Description Authenticate user with email/username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body dto.LoginRequest true "Login payload"
// @Success 200 {object} any
// @Failure 400 {object} any
// @Failure 401 {object} any
// @Router /api/v1/auth/login [post]
func (h *authHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.LoginRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			c.Error(err)
			return
		}

		loginResponse, err := h.authService.Login(c.Request.Context(), request)
		if err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = loginResponse

		c.JSON(http.StatusOK, res)
	}
}

// Register godoc
// @Summary Register
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body dto.RegisterRequest true "Register payload"
// @Success 201 {object} any
// @Failure 400 {object} any
// @Router /api/v1/auth/register [post]
func (h *authHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.RegisterRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			c.Error(err)
			return
		}

		userResponse, err := h.authService.Register(c.Request.Context(), request)
		if err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Status = dto.NewResponseStatus(http.StatusCreated)
		res.Request = dto.GetRequestClient(c)
		res.Data = userResponse

		c.JSON(http.StatusCreated, res)
	}
}

// RefreshToken godoc
// @Summary Refresh token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body dto.RefreshTokenRequest true "Refresh token payload"
// @Success 200 {object} any
// @Failure 400 {object} any
// @Router /api/v1/auth/refresh-token [post]
func (h *authHandler) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.RefreshTokenRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			c.Error(err)
			return
		}

		tokenResponse, err := h.authService.RefreshToken(c.Request.Context(), request)
		if err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = tokenResponse

		c.JSON(http.StatusOK, res)
	}
}

// ChangePassword godoc
// @Summary Change password
// @Description Change current user's password
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body dto.ResetPasswordRequest true "Change password payload"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} any
// @Failure 400 {object} any
// @Failure 401 {object} any
// @Router /api/v1/auth/change-password [put]
func (h *authHandler) ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.ChangePasswordRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			c.Error(err)
			return
		}
		id := c.GetString(consts.ContextUserID)
		if id == "" {
			c.Error(apperror.NewUnauthorizedError("User ID not found in context"))
			return
		}
		userID, err := uuid.Parse(id)
		if err != nil {
			c.Error(apperror.NewUnauthorizedError("Invalid user ID in context"))
			return
		}

		err = h.authService.ChangePassword(c, userID, request)
		if err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = gin.H{"message": "Password changed successfully", "user_id": userID}

		c.JSON(http.StatusOK, res)
	}
}

// ForgotPassword godoc
// @Summary Forgot password
// @Description Send reset password instructions to user's email
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body dto.ForgotPasswordRequest true "Forgot password payload"
// @Success 200 {object} any
// @Failure 400 {object} any
// @Router /api/v1/auth/forgot-password [post]
func (h *authHandler) ForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.ForgotPasswordRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			c.Error(err)
			return
		}

		err := h.authService.ForgotPassword(c.Request.Context(), request)
		if err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = gin.H{"message": "Reset password instructions sent to email if it exists"}

		c.JSON(http.StatusOK, res)
	}
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset user's password using reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body dto.ResetPasswordRequest true "Reset password payload"
// @Success 200 {object} any
// @Failure 400 {object} any
// @Router /api/v1/auth/reset-password [post]
func (h *authHandler) ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.ResetPasswordRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			c.Error(err)
			return
		}

		err := h.authService.ResetPassword(c.Request.Context(), request)
		if err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = gin.H{"message": "Password reset successfully"}

		c.JSON(http.StatusOK, res)
	}
}
