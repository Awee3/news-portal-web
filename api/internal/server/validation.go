package server

import (
	"errors"
	"strings"

	"news-portal/api/internal/database"
)

func validateLoginRequest(req *database.LoginRequest) error {
	if len(strings.TrimSpace(req.Email)) == 0 {
		return errors.New("email is required")
	}

	if len(req.Password) == 0 {
		return errors.New("password is required")
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	return nil
}

func validateUserRequest(req *database.UserRequest) error {
	if len(strings.TrimSpace(req.Username)) == 0 {
		return errors.New("username is required")
	}

	if len(req.Username) < 3 || len(req.Username) > 30 {
		return errors.New("username must be between 3 and 30 characters")
	}

	if len(strings.TrimSpace(req.Email)) == 0 {
		return errors.New("email is required")
	}

	if !isValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	if req.Role != "" && !isValidRole(req.Role) {
		return errors.New("role must be one of: admin, editor, user")
	}

	// Trim spaces and normalize
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	return nil
}

func validateUserUpdateRequest(req *database.UserUpdateRequest) error {
	if req.Username != "" {
		if len(req.Username) < 3 || len(req.Username) > 30 {
			return errors.New("username must be between 3 and 30 characters")
		}
		req.Username = strings.TrimSpace(req.Username)
	}

	if req.Email != "" {
		if !isValidEmail(req.Email) {
			return errors.New("invalid email format")
		}
		req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	}

	if req.Role != "" && !isValidRole(req.Role) {
		return errors.New("role must be one of: admin, editor, user")
	}

	return nil
}

func validatePasswordChangeRequest(req *database.PasswordChangeRequest) error {
	if len(req.CurrentPassword) == 0 {
		return errors.New("current password is required")
	}

	if len(req.NewPassword) < 6 {
		return errors.New("new password must be at least 6 characters")
	}

	if req.CurrentPassword == req.NewPassword {
		return errors.New("new password must be different from current password")
	}

	return nil
}

func isValidEmail(email string) bool {
	// Simple email validation
	return strings.Contains(email, "@") && strings.Contains(email, ".") && len(email) > 5
}

func isValidRole(role string) bool {
	validRoles := map[string]bool{
		"admin":  true,
		"editor": true,
		"user":   true,
	}
	return validRoles[role]
}
