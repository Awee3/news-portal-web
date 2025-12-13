package server

import (
	"context"

	"news-portal-web/api/internal/auth"
)

// Re-export context keys for use in handlers
const (
	UserIDKey   = auth.UserIDKey
	UserRoleKey = auth.UserRoleKey
	ClaimsKey   = auth.ClaimsKey
)

// Claims type alias for use in handlers
type Claims = auth.Claims

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	return auth.GetUserIDFromContext(ctx)
}

// GetUserRoleFromContext extracts user role from context
func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	return auth.GetUserRoleFromContext(ctx)
}

// GetClaimsFromContext extracts claims from context
func GetClaimsFromContext(ctx context.Context) (*Claims, bool) {
	return auth.GetClaimsFromContext(ctx)
}
