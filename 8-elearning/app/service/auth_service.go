package service

import (
	"context"
	"elearning-api/apperror"
	"elearning-api/config"
	"elearning-api/consts"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/repository"
	"elearning-api/util"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, request dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, request dto.RegisterRequest) (*dto.UserResponse, error)
	RefreshToken(ctx context.Context, request dto.RefreshTokenRequest) (*dto.TokenResponse, error)
	ChangePassword(ctx context.Context, request dto.ChangePasswordRequest) error
	ForgotPassword(ctx context.Context, request dto.ForgotPasswordRequest) error
}

type authService struct {
	userRepository repository.UserRepository
	jwtConfig      *config.JWTConfig
}

func NewAuthService(userRepository repository.UserRepository, jwtConfig *config.JWTConfig) AuthService {
	return &authService{
		userRepository: userRepository,
		jwtConfig:      jwtConfig,
	}
}

func (s *authService) Login(ctx context.Context, request dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get user by email or username
	user, err := s.userRepository.GetByEmailOrUsername(ctx, request.EmailOrUsername, request.EmailOrUsername)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get user", err)
	}

	if user == nil {
		return nil, apperror.NewUnauthorizedError("Username/email or password is incorrect")
	}

	// Compare password with hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, apperror.NewUnauthorizedError("Username/email or password is incorrect")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, apperror.NewForbiddenError("User account is inactive")
	}

	// Generate tokens
	accessToken, err := util.GenerateAccessToken(user.ID, string(user.Role), s.jwtConfig)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to generate access token", err)
	}

	refreshToken, err := util.GenerateRefreshToken(user.ID, s.jwtConfig)
	if err != nil {
		fmt.Println(err)
		return nil, apperror.NewInternalServerError("Failed to generate refresh token", err)
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
		return nil, apperror.NewInternalServerError("Failed to check existing user", err)
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
		return nil, apperror.NewInternalServerError("Failed to hash password", err)
	}

	newUser := &model.User{
		Email:    request.Email,
		Username: request.Username,
		Password: string(hashedPassword),
		Name:     request.Username,
		Role:     consts.RoleStudent,
	}

	createdUser, err := s.userRepository.Create(ctx, newUser)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create user", err)
	}

	return dto.NewUserDetailResponse(createdUser), nil
}

func (s *authService) RefreshToken(ctx context.Context, request dto.RefreshTokenRequest) (*dto.TokenResponse, error) {
	// Validate refresh token
	userID, err := util.ValidateRefreshToken(request.RefreshToken, s.jwtConfig)
	if err != nil {
		return nil, apperror.NewUnauthorizedError("Invalid or expired refresh token")
	}

	// Get user by ID
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get user", err)
	}

	if user == nil {
		return nil, apperror.NewUnauthorizedError("User not found")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, apperror.NewForbiddenError("User account is inactive")
	}

	// Generate new tokens
	accessToken, err := util.GenerateAccessToken(user.ID, string(user.Role), s.jwtConfig)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to generate access token", err)
	}

	newRefreshToken, err := util.GenerateRefreshToken(user.ID, s.jwtConfig)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to generate refresh token", err)
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *authService) ChangePassword(ctx context.Context, request dto.ChangePasswordRequest) error {
	// Get user ID from context
	userIDStr, exists := ctx.Value("user_id").(string)
	if !exists {
		return apperror.NewUnauthorizedError("User ID not found in context")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return apperror.NewUnauthorizedError("Invalid user ID in context")
	}

	// Get user by ID
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return apperror.NewInternalServerError("Failed to get user", err)
	}

	if user == nil {
		return apperror.NewUnauthorizedError("User not found")
	}

	// Compare old password with hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword))
	if err != nil {
		return apperror.NewUnauthorizedError("Old password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperror.NewInternalServerError("Failed to hash new password", err)
	}

	// Update user's password
	user.Password = string(hashedPassword)
	_, err = s.userRepository.Updates(ctx, user)
	if err != nil {
		return apperror.NewInternalServerError("Failed to update password", err)
	}

	return nil
}

func (s *authService) ForgotPassword(ctx context.Context, request dto.ForgotPasswordRequest) error {
	return nil
}
