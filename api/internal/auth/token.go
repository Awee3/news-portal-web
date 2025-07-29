package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"news-portal/api/internal/database"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type JWTManager struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey string, accessTokenTTL, refreshTokenTTL time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:       secretKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

// GenerateTokenPair generates both access and refresh tokens
func (j *JWTManager) GenerateTokenPair(user *database.User) (*TokenPair, error) {
	// Generate access token
	accessToken, err := j.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := j.generateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(j.accessTokenTTL.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// generateAccessToken creates a new access token
func (j *JWTManager) generateAccessToken(user *database.User) (string, error) {
	now := time.Now()

	claims := &Claims{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(user.UserID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "news-portal-api",
			Audience:  []string{"news-portal-web"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// generateRefreshToken creates a new refresh token
func (j *JWTManager) generateRefreshToken(user *database.User) (string, error) {
	now := time.Now()

	claims := &jwt.RegisteredClaims{
		Subject:   strconv.Itoa(user.UserID),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshTokenTTL)),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    "news-portal-api",
		Audience:  []string{"news-portal-refresh"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateAccessToken validates and parses access token
func (j *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ValidateRefreshToken validates refresh token
func (j *JWTManager) ValidateRefreshToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	// Check if it's a refresh token (audience check)
	if len(claims.Audience) == 0 || claims.Audience[0] != "news-portal-refresh" {
		return nil, errors.New("not a refresh token")
	}

	return claims, nil
}

// RefreshAccessToken generates new access token from refresh token
func (j *JWTManager) RefreshAccessToken(refreshToken string, user *database.User) (*TokenPair, error) {
	// Validate refresh token
	_, err := j.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Generate new token pair
	return j.GenerateTokenPair(user)
}

// ExtractUserIDFromToken extracts user ID from token without full validation
func (j *JWTManager) ExtractUserIDFromToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	return claims.UserID, nil
}

// IsTokenExpired checks if token is expired
func (j *JWTManager) IsTokenExpired(tokenString string) bool {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return true
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return true
	}

	return claims.ExpiresAt.Before(time.Now())
}

// RevokeToken adds token to blacklist (you might want to implement this with Redis/Database)
func (j *JWTManager) RevokeToken(tokenString string) error {
	// Implementation depends on your blacklist strategy
	// For now, just validate the token exists
	_, err := j.ValidateAccessToken(tokenString)
	if err != nil {
		return fmt.Errorf("cannot revoke invalid token: %w", err)
	}

	// TODO: Add to blacklist (Redis, Database, etc.)
	// For example:
	// return j.blacklistStore.Add(tokenString, claims.ExpiresAt.Time)

	return nil
}

// GetTokenClaims extracts all claims from token
func (j *JWTManager) GetTokenClaims(tokenString string) (*Claims, error) {
	return j.ValidateAccessToken(tokenString)
}
