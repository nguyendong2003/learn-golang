package handler

import (
	"elearning-api/dto"
	"elearning-api/service"
	"elearning-api/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Login() gin.HandlerFunc
	Register() gin.HandlerFunc
	RefreshToken() gin.HandlerFunc
	ChangePassword() gin.HandlerFunc
	ForgotPassword() gin.HandlerFunc
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

		c.JSON(http.StatusOK, dto.Success(loginResponse, "Login success"))
	}
}

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

		c.JSON(http.StatusCreated, dto.Success(userResponse, "Register success"))
	}
}

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

		c.JSON(http.StatusOK, dto.Success(tokenResponse, "Refresh token success"))
	}
}

func (h *authHandler) ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.ChangePasswordRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			c.Error(err)
			return
		}

		err := h.authService.ChangePassword(c.Request.Context(), request)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, dto.SimpleSuccess("Change password success"))
	}
}

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

		c.JSON(http.StatusOK, dto.SimpleSuccess("Password reset instructions sent"))
	}
}
