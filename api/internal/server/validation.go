package server

import (
	"errors"
	"regexp"
	"strings"

	"news-portal-web/api/internal/database"
)

// ========================================
// EMAIL VALIDATION
// ========================================

// isValidEmail validates email format
func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ========================================
// ROLE VALIDATION
// ========================================

// isValidRole validates user role
func isValidRole(role string) bool {
	return role == "admin" || role == "editor" || role == "user"
}

// ========================================
// STATUS VALIDATION
// ========================================

// isValidArticleStatus validates article status
func isValidArticleStatus(status string) bool {
	return status == "draft" || status == "published" || status == "archived"
}

// isValidCommentStatus validates comment status
func isValidCommentStatus(status string) bool {
	return status == "pending" || status == "approved" || status == "rejected"
}

// ========================================
// REQUEST VALIDATION
// ========================================

// validateLoginRequest validates login request
func validateLoginRequest(req *database.LoginRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}

	if !isValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

// validateUserRequest validates user registration request
func validateUserRequest(req *database.UserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}

	if len(req.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	if len(req.Username) > 50 {
		return errors.New("username must be at most 50 characters")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	if !isValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if req.Role != "" && !isValidRole(req.Role) {
		return errors.New("invalid role (must be admin, editor, or user)")
	}

	return nil
}

// validateUserUpdateRequest validates user update request
func validateUserUpdateRequest(req *database.UserUpdateRequest) error {
	if req.Username != "" && len(req.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	if req.Username != "" && len(req.Username) > 50 {
		return errors.New("username must be at most 50 characters")
	}

	if req.Email != "" && !isValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password != "" && len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if req.Role != "" && !isValidRole(req.Role) {
		return errors.New("invalid role (must be admin, editor, or user)")
	}

	return nil
}

// ========================================
// STRING HELPERS
// ========================================

// sanitizeString trims whitespace
func sanitizeString(s string) string {
	return strings.TrimSpace(s)
}

// truncateString truncates string to max length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
