package auth

import (
    "errors"
    "sync"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// TokenPair holds access and refresh tokens
type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"`
}

// Claims represents JWT claims
type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

// JWTManager handles JWT operations
type JWTManager struct {
    secretKey       []byte
    accessTokenTTL  time.Duration
    refreshTokenTTL time.Duration
    revokedTokens   map[string]time.Time
    mu              sync.RWMutex
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey string) *JWTManager {
    return &JWTManager{
        secretKey:       []byte(secretKey),
        accessTokenTTL:  60 * time.Minute,
        refreshTokenTTL: 7 * 24 * time.Hour,
        revokedTokens:   make(map[string]time.Time),
    }
}

// GenerateTokenPair generates access and refresh tokens
func (m *JWTManager) GenerateTokenPair(userID int, username, email, role string) (*TokenPair, error) {
    now := time.Now()

    // Access token
    accessClaims := Claims{
        UserID:   userID,
        Username: username,
        Email:    email,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   string(rune(userID)),
            IssuedAt:  jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTokenTTL)),
        },
    }

    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessTokenString, err := accessToken.SignedString(m.secretKey)
    if err != nil {
        return nil, err
    }

    // Refresh token
    refreshClaims := jwt.RegisteredClaims{
        Subject:   string(rune(userID)),
        IssuedAt:  jwt.NewNumericDate(now),
        ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTokenTTL)),
    }

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshTokenString, err := refreshToken.SignedString(m.secretKey)
    if err != nil {
        return nil, err
    }

    return &TokenPair{
        AccessToken:  accessTokenString,
        RefreshToken: refreshTokenString,
        ExpiresIn:    int64(m.accessTokenTTL.Seconds()),
    }, nil
}

// ValidateAccessToken validates an access token
func (m *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
    // Check if token is revoked
    m.mu.RLock()
    if _, revoked := m.revokedTokens[tokenString]; revoked {
        m.mu.RUnlock()
        return nil, errors.New("token has been revoked")
    }
    m.mu.RUnlock()

    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return m.secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}

// ValidateRefreshToken validates a refresh token
func (m *JWTManager) ValidateRefreshToken(tokenString string) (*jwt.RegisteredClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return m.secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid refresh token")
}

// RevokeToken adds a token to the revoked list
func (m *JWTManager) RevokeToken(tokenString string) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.revokedTokens[tokenString] = time.Now()
    return nil
}

// CleanupRevokedTokens removes expired tokens from the revoked list
func (m *JWTManager) CleanupRevokedTokens() {
    m.mu.Lock()
    defer m.mu.Unlock()

    now := time.Now()
    for token, revokedAt := range m.revokedTokens {
        // Remove tokens revoked more than refresh token TTL ago
        if now.Sub(revokedAt) > m.refreshTokenTTL {
            delete(m.revokedTokens, token)
        }
    }
}