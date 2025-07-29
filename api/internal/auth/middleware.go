package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
	ClaimsKey   contextKey = "claims"
)

// AuthMiddleware validates JWT token using JWTManager
func (j *JWTManager) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONError(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check Bearer format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			writeJSONError(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Validate token using JWTManager
		claims, err := j.ValidateAccessToken(token)
		if err != nil {
			writeJSONError(w, "Invalid or expired token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), UserIDKey, &claims.UserID)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
		ctx = context.WithValue(ctx, ClaimsKey, claims)

		// Call next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuthMiddleware validates JWT token if present, but doesn't require it
func (j *JWTManager) OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				if claims, err := j.ValidateAccessToken(token); err == nil {
					ctx := context.WithValue(r.Context(), UserIDKey, &claims.UserID)
					ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
					ctx = context.WithValue(ctx, ClaimsKey, claims)
					r = r.WithContext(ctx)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// Helper middleware functions (can be used without JWTManager instance)

// AdminOnly middleware - checks if user has admin role
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey).(*Claims)
		if !ok {
			writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims.Role != "admin" {
			writeJSONError(w, "Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// EditorOrAdmin middleware - for content management
func EditorOrAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey).(*Claims)
		if !ok {
			writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims.Role != "admin" && claims.Role != "editor" {
			writeJSONError(w, "Editor or Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireRole middleware - generic role checker
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ClaimsKey).(*Claims)
			if !ok {
				writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user role is in allowed roles
			allowed := false
			for _, role := range allowedRoles {
				if claims.Role == role {
					allowed = true
					break
				}
			}

			if !allowed {
				writeJSONError(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"error":"` + message + `"}`))
}
