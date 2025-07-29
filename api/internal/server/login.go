package server

import (
	"encoding/json"
	// "errors"
	"net/http"
	"strconv"
	"strings"

	"news-portal/api/internal/auth"
	"news-portal/api/internal/database"

	"github.com/gorilla/mux"
)

func (s *Server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req database.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := validateLoginRequest(&req); err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := database.AuthenticateUser(r.Context(), s.GetDB(), &req)
		if err != nil {
			writeJSONError(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Generate JWT tokens
		tokenPair, err := s.jwtManager.GenerateTokenPair(user)
		if err != nil {
			writeJSONError(w, "Failed to generate tokens", http.StatusInternalServerError)
			return
		}

		// User profile without password
		profile := database.UserProfile{
			UserID:    user.UserID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		}

		response := map[string]interface{}{
			"message": "Login successful",
			"user":    profile,
			"tokens":  tokenPair,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func (s *Server) handleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req database.UserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := validateUserRequest(&req); err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Set default role to "user" for registration
		if req.Role == "" {
			req.Role = "user"
		}

		// Check if email already exists
		emailExists, err := database.IsEmailExists(r.Context(), s.GetDB(), req.Email)
		if err != nil {
			writeJSONError(w, "Failed to check email existence: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if emailExists {
			writeJSONError(w, "Email already exists", http.StatusConflict)
			return
		}

		// Check if username already exists
		usernameExists, err := database.IsUsernameExists(r.Context(), s.GetDB(), req.Username)
		if err != nil {
			writeJSONError(w, "Failed to check username existence: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if usernameExists {
			writeJSONError(w, "Username already exists", http.StatusConflict)
			return
		}

		user, err := database.CreateUser(r.Context(), s.GetDB(), &req)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				writeJSONError(w, "Email or username already exists", http.StatusConflict)
				return
			}
			writeJSONError(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Generate JWT tokens for new user
		tokenPair, err := s.jwtManager.GenerateTokenPair(user)
		if err != nil {
			writeJSONError(w, "User created but failed to generate tokens", http.StatusInternalServerError)
			return
		}

		// User profile without password
		profile := database.UserProfile{
			UserID:    user.UserID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		}

		response := map[string]interface{}{
			"message": "Registration successful",
			"user":    profile,
			"tokens":  tokenPair,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func (s *Server) handleRefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if req.RefreshToken == "" {
			writeJSONError(w, "Refresh token is required", http.StatusBadRequest)
			return
		}

		// Validate refresh token and get user ID
		claims, err := s.jwtManager.ValidateRefreshToken(req.RefreshToken)
		if err != nil {
			writeJSONError(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}

		// Get user from database
		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			writeJSONError(w, "Invalid token claims", http.StatusBadRequest)
			return
		}

		user, err := database.GetUserByID(r.Context(), s.GetDB(), userID)
		if err != nil {
			writeJSONError(w, "User not found", http.StatusNotFound)
			return
		}

		// Generate new token pair
		tokenPair, err := s.jwtManager.RefreshAccessToken(req.RefreshToken, user)
		if err != nil {
			writeJSONError(w, "Failed to refresh token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokenPair)
	}
}

func (s *Server) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from header to revoke it
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONError(w, "No token to revoke", http.StatusBadRequest)
			return
		}

		if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			writeJSONError(w, "Invalid authorization header format", http.StatusBadRequest)
			return
		}

		tokenString := authHeader[7:] // Remove "Bearer " prefix

		// Revoke token (add to blacklist)
		if err := s.jwtManager.RevokeToken(tokenString); err != nil {
			writeJSONError(w, "Failed to revoke token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Logged out successfully",
		})
	}
}

func (s *Server) handleGetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context (set by auth middleware)
		userID, ok := r.Context().Value(auth.UserIDKey).(*int)
		if !ok {
			writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := database.GetUserByID(r.Context(), s.GetDB(), *userID)
		if err != nil {
			writeJSONError(w, "User not found", http.StatusNotFound)
			return
		}

		profile := database.UserProfile{
			UserID:    user.UserID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	}
}

func (s *Server) handleUpdateProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID, ok := r.Context().Value(auth.UserIDKey).(*int)
		if !ok {
			writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req database.UserUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := validateUserUpdateRequest(&req); err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if new email already exists (excluding current user)
		if req.Email != "" {
			existing, err := database.GetUserByEmail(r.Context(), s.GetDB(), req.Email)
			if err == nil && existing.UserID != *userID {
				writeJSONError(w, "Email already exists", http.StatusConflict)
				return
			}
		}

		// Check if new username already exists (excluding current user)
		if req.Username != "" {
			existing, err := database.GetUserByUsername(r.Context(), s.GetDB(), req.Username)
			if err == nil && existing.UserID != *userID {
				writeJSONError(w, "Username already exists", http.StatusConflict)
				return
			}
		}

		// Users can't change their own role
		req.Role = ""

		user, err := database.UpdateUser(r.Context(), s.GetDB(), *userID, &req)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				writeJSONError(w, "User not found", http.StatusNotFound)
				return
			}
			if strings.Contains(err.Error(), "duplicate") {
				writeJSONError(w, "Email or username already exists", http.StatusConflict)
				return
			}
			writeJSONError(w, "Failed to update profile: "+err.Error(), http.StatusInternalServerError)
			return
		}

		profile := database.UserProfile{
			UserID:    user.UserID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(profile)
	}
}

func (s *Server) handleChangePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID, ok := r.Context().Value(auth.UserIDKey).(*int)
		if !ok {
			writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req database.PasswordChangeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := validatePasswordChangeRequest(&req); err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := database.ChangePassword(r.Context(), s.GetDB(), *userID, &req); err != nil {
			if strings.Contains(err.Error(), "not found") {
				writeJSONError(w, "User not found", http.StatusNotFound)
				return
			}
			if strings.Contains(err.Error(), "incorrect") {
				writeJSONError(w, "Current password is incorrect", http.StatusBadRequest)
				return
			}
			writeJSONError(w, "Failed to change password: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Password changed successfully",
		})
	}
}

func (s *Server) RegisterAuthRoutes(router *mux.Router) {
    authRouter := router.PathPrefix("/auth").Subrouter()

    // Public auth routes
    authRouter.HandleFunc("/register", s.handleRegister()).Methods("POST")
    authRouter.HandleFunc("/login", s.handleLogin()).Methods("POST")
    authRouter.HandleFunc("/refresh", s.handleRefreshToken()).Methods("POST")

    // Protected auth routes
    protectedAuthRouter := authRouter.PathPrefix("").Subrouter()
    protectedAuthRouter.Use(s.jwtManager.AuthMiddleware)
    protectedAuthRouter.HandleFunc("/logout", s.handleLogout()).Methods("POST")
    protectedAuthRouter.HandleFunc("/profile", s.handleGetProfile()).Methods("GET")
    protectedAuthRouter.HandleFunc("/profile", s.handleUpdateProfile()).Methods("PUT")
    protectedAuthRouter.HandleFunc("/change-password", s.handleChangePassword()).Methods("POST")
}