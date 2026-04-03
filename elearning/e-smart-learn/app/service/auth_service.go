package service

import (
	"context"
	"fmt"
	"strings"
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
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

type AuthService interface {
	Login(ctx context.Context, request dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, request dto.RegisterRequest) (*dto.UserResponse, error)
	GetGoogleLoginURL(ctx context.Context) (string, error)
	GoogleCallback(ctx context.Context, state, code string) (*dto.LoginResponse, error)
	LoginWithGoogle(ctx context.Context, payload dto.GoogleOAuthPayload) (*dto.LoginResponse, error)
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
	googleOAuthConfig      *config.GoogleOAuthConfig
	googleOAuth            *oauth2.Config
	mail                   pkg.EmailProvider
	cache                  pkg.CacheProvider
}

func NewAuthService(
	userRepository repository.UserRepository,
	roleRepository repository.RoleRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	jwtConfig *config.JWTConfig,
	googleOAuthConfig *config.GoogleOAuthConfig,
	mail pkg.EmailProvider,
	cache pkg.CacheProvider,
) AuthService {
	var oauth2Cfg *oauth2.Config
	if googleOAuthConfig != nil {
		oauth2Cfg = googleOAuthConfig.OAuth2Config()
	}

	return &authService{
		userRepository:         userRepository,
		roleRepository:         roleRepository,
		refreshTokenRepository: refreshTokenRepository,
		jwtConfig:              jwtConfig,
		googleOAuthConfig:      googleOAuthConfig,
		googleOAuth:            oauth2Cfg,
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
		Email:         request.Email,
		Username:      request.Username,
		Password:      string(hashedPassword),
		Name:          request.Username,
		OauthProvider: nil,
		OauthID:       nil,
		Role:          role,
		Avatar:        fmt.Sprintf("https://i.pravatar.cc/150?u=%s", request.Username),
	}

	createdUser, err := s.userRepository.Create(ctx, newUser)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create user")
	}

	return dto.NewUserDetailResponse(createdUser), nil
}

func (s *authService) GetGoogleLoginURL(ctx context.Context) (string, error) {
	if s.googleOAuth == nil {
		return "", apperror.NewInternalServerError("Google OAuth is not configured")
	}

	state := util.GenerateOAuthState()
	if state == "" {
		return "", apperror.NewInternalServerError("Failed to generate OAuth state")
	}

	cacheKey := fmt.Sprintf("oauth_google_state_%s", state)
	if err := s.cache.Set(ctx, cacheKey, "1", 10*time.Minute); err != nil {
		return "", apperror.NewInternalServerError("Failed to save OAuth state")
	}

	return s.googleOAuth.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (s *authService) GoogleCallback(ctx context.Context, state, code string) (*dto.LoginResponse, error) {
	if s.googleOAuth == nil || s.googleOAuthConfig == nil {
		return nil, apperror.NewInternalServerError("Google OAuth is not configured")
	}

	if state == "" {
		return nil, apperror.NewBadRequestError("Missing OAuth state")
	}

	if code == "" {
		return nil, apperror.NewBadRequestError("Missing authorization code")
	}

	cacheKey := fmt.Sprintf("oauth_google_state_%s", state)
	if _, err := s.cache.Get(ctx, cacheKey); err != nil {
		return nil, apperror.NewBadRequestError("Invalid OAuth state")
	}
	_ = s.cache.Delete(ctx, cacheKey)

	token, err := s.googleOAuth.Exchange(ctx, code)
	if err != nil {
		return nil, apperror.NewUnauthorizedError("Failed to exchange authorization code")
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok || rawIDToken == "" {
		return nil, apperror.NewUnauthorizedError("Missing id_token from Google response")
	}

	payload, err := idtoken.Validate(ctx, rawIDToken, s.googleOAuthConfig.ClientID)
	if err != nil {
		return nil, apperror.NewUnauthorizedError("Invalid Google id_token")
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)
	sub, _ := payload.Claims["sub"].(string)
	emailVerified, hasEmailVerified := payload.Claims["email_verified"].(bool)

	if email == "" || sub == "" {
		return nil, apperror.NewUnauthorizedError("Invalid Google account payload")
	}

	if hasEmailVerified && !emailVerified {
		return nil, apperror.NewUnauthorizedError("Google email is not verified")
	}

	return s.LoginWithGoogle(ctx, dto.GoogleOAuthPayload{
		Email:   email,
		Name:    name,
		Picture: picture,
		Sub:     sub,
	})
}

func (s *authService) LoginWithGoogle(ctx context.Context, payload dto.GoogleOAuthPayload) (*dto.LoginResponse, error) {
	if payload.Email == "" || payload.Sub == "" {
		return nil, apperror.NewBadRequestError("Invalid Google payload")
	}

	user, err := s.userRepository.GetByOAuth(ctx, "google", payload.Sub)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get user by OAuth")
	}

	if user == nil {
		user, err = s.userRepository.GetByEmail(ctx, payload.Email)
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to get user by email")
		}
	}

	if user == nil {
		role, err := s.roleRepository.GetByName(ctx, string(consts.RoleStudent))
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to get student role")
		}
		if role == nil {
			return nil, apperror.NewNotFoundError("Role student not found")
		}

		username := s.generateOAuthUsername(payload.Email)
		for range 5 {
			exists, err := s.userRepository.CheckExists(ctx, "username = ?", username)
			if err != nil {
				return nil, apperror.NewInternalServerError("Failed to validate username")
			}
			if !exists {
				break
			}
			username = fmt.Sprintf("%s_%s", username, util.GenerateRandomNumber(4))
		}

		randomPassword := util.GenerateRandomString(16)
		if randomPassword == "" {
			randomPassword = fmt.Sprintf("oauth_%d", time.Now().UnixNano())
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(randomPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to generate password")
		}

		name := payload.Name
		if name == "" {
			name = username
		}

		provider := "google"

		newUser := &model.User{
			Email:            payload.Email,
			Username:         username,
			Password:         string(hashedPassword),
			Name:             name,
			Avatar:           payload.Picture,
			OauthProvider:    &provider,
			OauthID:          &payload.Sub,
			IsActive:         true,
			RoleID:           role.ID,
			StripeCustomerID: nil,
		}

		created, err := s.userRepository.Create(ctx, newUser)
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to create Google user")
		}

		user, err = s.userRepository.FindByID(ctx, created.ID, []repository.Preload{repository.PreloadPath(repository.Role)})
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to load created user")
		}
	} else {
		if user.OauthID != nil && *user.OauthID != payload.Sub {
			return nil, apperror.NewUnauthorizedError("This email is linked to another Google account")
		}

		updates := map[string]any{}
		if user.OauthID == nil {
			updates["oauth_provider"] = "google"
			updates["oauth_id"] = &payload.Sub
		}
		if user.Avatar == "" && payload.Picture != "" {
			updates["avatar"] = payload.Picture
		}
		if user.Name == "" && payload.Name != "" {
			updates["name"] = payload.Name
		}

		if len(updates) > 0 {
			_, err := s.userRepository.Update(ctx, user.ID, updates)
			if err != nil {
				return nil, apperror.NewInternalServerError("Failed to link Google account")
			}
		}

		user, err = s.userRepository.FindByID(ctx, user.ID, []repository.Preload{repository.PreloadPath(repository.Role)})
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to load user")
		}
	}

	if user == nil {
		return nil, apperror.NewInternalServerError("Failed to resolve user")
	}

	if !user.IsActive {
		return nil, apperror.NewForbiddenError("User account is inactive")
	}

	accessToken, err := util.GenerateAccessToken(user.ID, s.jwtConfig)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to generate access token")
	}

	refreshToken, expirationTime, err := util.GenerateRefreshToken(user.ID, s.jwtConfig)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to generate refresh token")
	}

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

func (s *authService) generateOAuthUsername(email string) string {
	base := strings.ToLower(strings.TrimSpace(email))
	if at := strings.Index(base, "@"); at > 0 {
		base = base[:at]
	}
	base = strings.ReplaceAll(base, " ", "")
	base = strings.ReplaceAll(base, ".", "_")
	base = strings.ReplaceAll(base, "-", "_")
	if base == "" {
		base = "google_user"
	}
	return base
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
	log := util.LoggerFromContextWithLayer(ctx, util.LayerService)
	if err := s.cache.Delete(ctx, fmt.Sprintf("password_reset_%s", request.Email)); err != nil {
		log.Warn("Failed to clear reset code from cache", "email", request.Email, "error", err)
	}

	return nil
}
