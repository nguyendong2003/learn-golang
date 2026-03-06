package util

import (
	"elearning-api/config"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	AccessTokenType  TokenType = "access_token"
	RefreshTokenType TokenType = "refresh_token"
)

type AccessTokenClaims struct {
	UserID string    `json:"user_id"`
	Role   string    `json:"role"`
	Type   TokenType `json:"type"`

	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID string    `json:"user_id"`
	Type   TokenType `json:"type"`

	jwt.RegisteredClaims
}

// GenerateAccessToken generates a new JWT access token for a user
func GenerateAccessToken(userID uuid.UUID, role string, jwtConfig *config.JWTConfig) (string, error) {
	expirationTime := time.Now().Add(jwtConfig.AccessTokenExpiration)

	claims := &AccessTokenClaims{
		UserID: userID.String(),
		Role:   role,
		Type:   AccessTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtConfig.Issuer,
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtConfig.AccessTokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a new JWT refresh token for a user
func GenerateRefreshToken(userID uuid.UUID, jwtConfig *config.JWTConfig) (string, error) {
	expirationTime := time.Now().Add(jwtConfig.RefreshTokenExpiration)

	claims := &RefreshTokenClaims{
		UserID: userID.String(),
		Type:   RefreshTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtConfig.Issuer,
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtConfig.RefreshTokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateAccessToken validates the JWT access token and returns the claims
func ValidateAccessToken(tokenString string, jwtConfig *config.JWTConfig) (*AccessTokenClaims, error) {
	jwtAccessTokenSecret := jwtConfig.AccessTokenSecret
	jwtIssuer := jwtConfig.Issuer

	token, err := jwt.ParseWithClaims(
		tokenString,
		&AccessTokenClaims{},
		func(token *jwt.Token) (any, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}

			return []byte(jwtAccessTokenSecret), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("access token expired")
		}
		return nil, err
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	if claims.Type != AccessTokenType {
		return nil, errors.New("invalid token type")
	}

	if claims.Issuer != jwtIssuer {
		return nil, errors.New("invalid token issuer")
	}

	return claims, nil
}

// ValidateRefreshToken validates the JWT refresh token and returns the user ID
func ValidateRefreshToken(tokenString string, jwtConfig *config.JWTConfig) (uuid.UUID, error) {
	jwtRefreshTokenSecret := jwtConfig.RefreshTokenSecret
	jwtIssuer := jwtConfig.Issuer
	token, err := jwt.ParseWithClaims(
		tokenString,
		&RefreshTokenClaims{},
		func(token *jwt.Token) (any, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}

			return []byte(jwtRefreshTokenSecret), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return uuid.Nil, errors.New("refresh token expired")
		}
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid refresh token")
	}

	if claims.Type != RefreshTokenType {
		return uuid.Nil, errors.New("invalid token type")
	}

	if claims.Issuer != jwtIssuer {
		return uuid.Nil, errors.New("invalid token issuer")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in token: %w", err)
	}

	return userID, nil
}

// ExtractTokenFromHeader extracts the JWT token from Authorization header
// Expected format: "Bearer <token>"
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is empty")
	}

	// Check if the header starts with "Bearer "
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) {
		return "", errors.New("invalid authorization header format")
	}

	if authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("authorization header must start with 'Bearer '")
	}

	token := authHeader[len(bearerPrefix):]
	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil
}
