package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"news-portal-web/api/internal/auth"
	"news-portal-web/api/internal/database"

	"github.com/gorilla/mux"
)

// ========================================
// AUTH HANDLERS
// ========================================

// handleLogin - POST /api/v1/auth/login
func (s *Server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req database.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := validateLoginRequest(&req); err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := database.AuthenticateUser(r.Context(), s.GetDB(), &req)
		if err != nil {
			writeJSONError(w, "Email atau password salah", http.StatusUnauthorized)
			return
		}

		// Generate JWT tokens
		jwtManager := s.GetJWTManager()
		if jwtManager == nil {
			writeJSONError(w, "JWT manager not configured", http.StatusInternalServerError)
			return
		}

		tokenPair, err := jwtManager.GenerateTokenPair(user.UserID, user.Username, user.Email, user.Role)
		if err != nil {
			writeJSONError(w, "Failed to generate tokens", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"message": "Login berhasil",
			"user":    user.ToPublic(),
			"tokens":  tokenPair,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// handleRegister - POST /api/v1/auth/register
func (s *Server) handleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req database.UserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := validateUserRequest(&req); err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Default role
		if req.Role == "" {
			req.Role = "user"
		}

		// Check email exists
		emailExists, err := database.IsEmailExists(r.Context(), s.GetDB(), req.Email)
		if err != nil {
			writeJSONError(w, "Error checking email", http.StatusInternalServerError)
			return
		}
		if emailExists {
			writeJSONError(w, "Email sudah terdaftar", http.StatusConflict)
			return
		}

		// Check username exists
		usernameExists, err := database.IsUsernameExists(r.Context(), s.GetDB(), req.Username)
		if err != nil {
			writeJSONError(w, "Error checking username", http.StatusInternalServerError)
			return
		}
		if usernameExists {
			writeJSONError(w, "Username sudah digunakan", http.StatusConflict)
			return
		}

		user, err := database.CreateUser(r.Context(), s.GetDB(), &req)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				writeJSONError(w, "Email atau username sudah terdaftar", http.StatusConflict)
				return
			}
			writeJSONError(w, "Gagal membuat user", http.StatusInternalServerError)
			return
		}

		// Generate JWT tokens
		jwtManager := s.GetJWTManager()
		if jwtManager == nil {
			writeJSONError(w, "JWT manager not configured", http.StatusInternalServerError)
			return
		}

		tokenPair, err := jwtManager.GenerateTokenPair(user.UserID, user.Username, user.Email, user.Role)
		if err != nil {
			writeJSONError(w, "Failed to generate tokens", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"message": "Registrasi berhasil",
			"user":    user.ToPublic(),
			"tokens":  tokenPair,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// handleRefreshToken - POST /api/v1/auth/refresh
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
			writeJSONError(w, "Refresh token harus diisi", http.StatusBadRequest)
			return
		}

		jwtManager := s.GetJWTManager()
		if jwtManager == nil {
			writeJSONError(w, "JWT manager not configured", http.StatusInternalServerError)
			return
		}

		// Validate refresh token
		claims, err := jwtManager.ValidateRefreshToken(req.RefreshToken)
		if err != nil {
			writeJSONError(w, "Refresh token tidak valid", http.StatusUnauthorized)
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
			writeJSONError(w, "User tidak ditemukan", http.StatusNotFound)
			return
		}

		// Generate new token pair
		tokenPair, err := jwtManager.GenerateTokenPair(user.UserID, user.Username, user.Email, user.Role)
		if err != nil {
			writeJSONError(w, "Failed to refresh token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokenPair)
	}
}

// handleLogout - POST /api/v1/auth/logout
func (s *Server) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			writeJSONError(w, "Token tidak ditemukan", http.StatusBadRequest)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		jwtManager := s.GetJWTManager()
		if jwtManager == nil {
			writeJSONError(w, "JWT manager not configured", http.StatusInternalServerError)
			return
		}

		// Revoke token
		if err := jwtManager.RevokeToken(tokenString); err != nil {
			writeJSONError(w, "Gagal logout", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Logout berhasil",
		})
	}
}

func (s *Server) handleGetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context (set by auth middleware)
		userID, ok := r.Context().Value(auth.UserIDKey).(int)
		if !ok {
			// Try pointer version
			userIDPtr, ok := r.Context().Value(auth.UserIDKey).(*int)
			if !ok || userIDPtr == nil {
				writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			userID = *userIDPtr
		}

		user, err := database.GetUserByID(r.Context(), s.GetDB(), userID)
		if err != nil {
			writeJSONError(w, "User not found", http.StatusNotFound)
			return
		}

		profile := user.ToProfile()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	}
}

func (s *Server) handleUpdateProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID, ok := r.Context().Value(auth.UserIDKey).(int)
		if !ok {
			userIDPtr, ok := r.Context().Value(auth.UserIDKey).(*int)
			if !ok || userIDPtr == nil {
				writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			userID = *userIDPtr
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
			if err == nil && existing.UserID != userID {
				writeJSONError(w, "Email already exists", http.StatusConflict)
				return
			}
		}

		// Check if new username already exists (excluding current user)
		if req.Username != "" {
			existing, err := database.GetUserByUsername(r.Context(), s.GetDB(), req.Username)
			if err == nil && existing.UserID != userID {
				writeJSONError(w, "Username already exists", http.StatusConflict)
				return
			}
		}

		// Users can't change their own role
		req.Role = ""

		user, err := database.UpdateUser(r.Context(), s.GetDB(), userID, &req)
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

		profile := user.ToProfile()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(profile)
	}
}

// ========================================
// CHANGE PASSWORD HANDLER
// ========================================

// ChangePasswordRequest - Request body untuk ganti password

// ========================================
// ROUTE REGISTRATION
// ========================================

func (s *Server) RegisterAuthRoutes(router *mux.Router) {
	authRouter := router.PathPrefix("/auth").Subrouter()

	// Public auth routes
	authRouter.HandleFunc("/register", s.handleRegister()).Methods("POST")
	authRouter.HandleFunc("/login", s.handleLogin()).Methods("POST")
	authRouter.HandleFunc("/refresh", s.handleRefreshToken()).Methods("POST")

	// Protected auth routes
	jwtManager := s.GetJWTManager()
	if jwtManager != nil {
		protectedAuthRouter := authRouter.PathPrefix("").Subrouter()
		protectedAuthRouter.Use(auth.AuthMiddleware(jwtManager))
		protectedAuthRouter.HandleFunc("/logout", s.handleLogout()).Methods("POST")
		protectedAuthRouter.HandleFunc("/profile", s.handleGetProfile()).Methods("GET")
		protectedAuthRouter.HandleFunc("/profile", s.handleUpdateProfile()).Methods("PUT")
		protectedAuthRouter.HandleFunc("/change-password", s.handleChangePassword()).Methods("POST")
	}
}
