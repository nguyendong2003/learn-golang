package service

import (
	"context"
	"fmt"
	"time"

	"elearning-api/apperror"
	"elearning-api/config"
	"elearning-api/consts"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/pkg"
	"elearning-api/repository"
	"elearning-api/util"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, request dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, request dto.RegisterRequest) (*dto.UserResponse, error)
	RefreshToken(ctx context.Context, request dto.RefreshTokenRequest) (*dto.TokenResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, request dto.ChangePasswordRequest) error
	ForgotPassword(ctx context.Context, request dto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, request dto.ResetPasswordRequest) error
}

type authService struct {
	userRepository         repository.UserRepository
	roleRepository         repository.RoleRepository
	refreshTokenRepository repository.RefreshTokenRepository
	jwtConfig              *config.JWTConfig
	mail                   pkg.EmailProvider
	cache                  pkg.CacheProvider
}

func NewAuthService(
	userRepository repository.UserRepository,
	roleRepository repository.RoleRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	jwtConfig *config.JWTConfig,
	mail pkg.EmailProvider,
	cache pkg.CacheProvider,
) AuthService {
	return &authService{
		userRepository:         userRepository,
		roleRepository:         roleRepository,
		refreshTokenRepository: refreshTokenRepository,
		jwtConfig:              jwtConfig,
		mail:                   mail,
		cache:                  cache,
	}
}

func (s *authService) Login(ctx context.Context, request dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get user by email or username
	user, err := s.userRepository.GetByEmailOrUsername(ctx, request.EmailOrUsername, request.EmailOrUsername)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get user")
	}

	if user == nil {
		return nil, apperror.NewUnauthorizedError("Username or password is incorrect")
	}

	// Compare password with hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, apperror.NewUnauthorizedError("Username or password is incorrect")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, apperror.NewForbiddenError("User account is inactive")
	}

	// Generate tokens
	accessToken, err := util.GenerateAccessToken(user.ID, s.jwtConfig)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to generate access token")
	}

	refreshToken, expirationTime, err := util.GenerateRefreshToken(user.ID, s.jwtConfig)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to generate refresh token")
	}

	// Save refresh token
	refreshTokenModel := &model.RefreshToken{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiredAt: expirationTime,
	}

	_, err = s.refreshTokenRepository.Create(ctx, refreshTokenModel)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to save refresh token")
	}

	return &dto.LoginResponse{
		TokenResponse: dto.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		User: dto.NewUserDetailResponse(user),
	}, nil
}

func (s *authService) Register(ctx context.Context, request dto.RegisterRequest) (*dto.UserResponse, error) {
	user, err := s.userRepository.GetByEmailOrUsername(ctx, request.Email, request.Username)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to check existing user")
	}

	if user != nil {
		validationErrors := map[string]string{}

		if user.Email == request.Email {
			validationErrors["email"] = "Email already exists"
		}

		if user.Username == request.Username {
			validationErrors["username"] = "Username already exists"
		}

		return nil, apperror.NewValidationError(validationErrors)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to hash password")
	}

	role, err := s.roleRepository.GetByName(ctx, string(consts.RoleStudent))
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get role")
	}

	if role == nil {
		return nil, apperror.NewNotFoundError("Role not found")
	}

	newUser := &model.User{
		Email:    request.Email,
		Username: request.Username,
		Password: string(hashedPassword),
		Name:     request.Username,
		Role:     role,
	}

	createdUser, err := s.userRepository.Create(ctx, newUser)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create user")
	}

	return dto.NewUserDetailResponse(createdUser), nil
}

func (s *authService) RefreshToken(ctx context.Context, request dto.RefreshTokenRequest) (*dto.TokenResponse, error) {
	// Validate refresh token
	_, err := util.ValidateRefreshToken(request.RefreshToken, s.jwtConfig)
	if err != nil {
		return nil, apperror.NewUnauthorizedError("Invalid or expired refresh token")
	}

	refreshToken, err := s.refreshTokenRepository.GetByToken(ctx, request.RefreshToken)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get refresh token")
	}

	if refreshToken == nil {
		return nil, apperror.NewNotFoundError("Refresh token not found")
	}

	if refreshToken.IsRevoked || refreshToken.ExpiredAt.Before(time.Now()) {
		return nil, apperror.NewUnauthorizedError("Invalid or expired refresh token")
	}

	// Get user by ID
	user := refreshToken.User
	if user == nil {
		return nil, apperror.NewNotFoundError("User not found")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, apperror.NewForbiddenError("User account is inactive")
	}

	// Generate new tokens
	accessToken, err := util.GenerateAccessToken(user.ID, s.jwtConfig)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to generate access token")
	}

	newRefreshToken, expirationTime, err := util.GenerateRefreshToken(user.ID, s.jwtConfig)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to generate refresh token")
	}

	// Create new refresh token model
	newRefreshTokenModel := &model.RefreshToken{
		Token:     newRefreshToken,
		UserID:    user.ID,
		ExpiredAt: expirationTime,
	}

	_, err = s.refreshTokenRepository.Create(ctx, newRefreshTokenModel)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to save new refresh token")
	}

	// Delete old refresh token
	err = s.refreshTokenRepository.Delete(ctx, refreshToken.ID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to delete old refresh token")
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *authService) ChangePassword(ctx context.Context, userID uuid.UUID, request dto.ChangePasswordRequest) error {
	if request.NewPassword != request.ConfirmPassword {
		return apperror.NewBadRequestError("New password and confirm password do not match")
	}

	// Get user by ID
	user, err := s.userRepository.FindByID(ctx, userID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to get user")
	}
	if user == nil {
		return apperror.NewUnauthorizedError("User not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword))
	if err != nil {
		return apperror.NewUnauthorizedError("Old password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperror.NewInternalServerError("Failed to hash new password")
	}

	user.Password = string(hashedPassword)
	_, err = s.userRepository.Updates(ctx, user)
	if err != nil {
		return apperror.NewInternalServerError("Failed to update password")
	}

	return nil
}

func (s *authService) ForgotPassword(ctx context.Context, request dto.ForgotPasswordRequest) error {
	user, err := s.userRepository.GetByEmail(ctx, request.Email)
	if err != nil {
		return apperror.NewInternalServerError("Failed to get user")
	}
	if user == nil {
		return apperror.NewBadRequestError("Invalid email address")
	}
	resetCode := util.GenerateRandomNumber(6)
	cacheKey := fmt.Sprintf("password_reset_%s", request.Email)
	err = s.cache.Set(ctx, cacheKey, resetCode, 15*time.Minute)
	if err != nil {
		return apperror.NewInternalServerError("Failed to set reset code in cache")
	}
	err = s.mail.SendPasswordReset(user.Email, user.Name, resetCode)
	if err != nil {
		return apperror.NewInternalServerError("Failed to send password reset email")
	}
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, request dto.ResetPasswordRequest) error {
	resetCode, err := s.cache.Get(ctx, fmt.Sprintf("password_reset_%s", request.Email))
	if err != nil {
		return apperror.NewBadRequestError("Invalid or expired reset code")
	}
	if resetCode != request.ResetCode {
		return apperror.NewBadRequestError("Invalid reset code")
	}
	if request.NewPassword != request.ConfirmPassword {
		return apperror.NewBadRequestError("New password and confirm password do not match")
	}
	user, err := s.userRepository.GetByEmail(ctx, request.Email)
	if err != nil {
		return apperror.NewInternalServerError("Failed to get user")
	}
	if user == nil {
		return apperror.NewBadRequestError("Invalid email address")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperror.NewInternalServerError("Failed to hash new password")
	}

	user.Password = string(hashedPassword)
	_, err = s.userRepository.Updates(ctx, user)
	if err != nil {
		return apperror.NewInternalServerError("Failed to update password")
	}
	s.cache.Delete(ctx, fmt.Sprintf("password_reset_%s", request.Email))

	return nil
}
