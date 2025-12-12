package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"news-portal/api/internal/auth"
	"news-portal/api/internal/database"

	"github.com/gorilla/mux"
)
// Admin-only user management functions
func (s *Server) handleCreateUserByAdmin() http.HandlerFunc {
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

		// Convert to profile (without password)
		profile := database.UserProfile{
			UserID:    user.UserID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(profile)
	}
}

func (s *Server) handleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userIDStr := vars["id"]

		if userIDStr != "" {
			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				writeJSONError(w, "Invalid user ID", http.StatusBadRequest)
				return
			}

			user, err := database.GetUserByID(r.Context(), s.GetDB(), userID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found") {
					writeJSONError(w, "User not found", http.StatusNotFound)
				} else {
					writeJSONError(w, "Failed to get user: "+err.Error(), http.StatusInternalServerError)
				}
				return
			}

			// Convert to profile (without password)
			profile := database.UserProfile{
				UserID:    user.UserID,
				Username:  user.Username,
				Email:     user.Email,
				Role:      user.Role,
				CreatedAt: user.CreatedAt,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(profile)
			return
		}

		// List users
		s.handleListUsers()(w, r)
	}
}

func (s *Server) handleListUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		role := r.URL.Query().Get("role")
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		// Filter by role
		if role != "" {
			users, err := database.ListUsersByRole(r.Context(), s.GetDB(), role)
			if err != nil {
				writeJSONError(w, "Failed to fetch users: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"users": users,
			})
			return
		}

		// List all users with pagination
		limit := 10
		offset := 0

		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}

		if offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
				offset = o
			}
		}

		users, totalCount, err := database.ListUsers(r.Context(), s.GetDB(), limit, offset)
		if err != nil {
			writeJSONError(w, "Failed to fetch users: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Calculate pagination info
		totalPages := (totalCount + limit - 1) / limit
		currentPage := (offset / limit) + 1

		response := map[string]interface{}{
			"users": users,
			"pagination": map[string]interface{}{
				"current_page": currentPage,
				"total_pages":  totalPages,
				"total_count":  totalCount,
				"limit":        limit,
				"offset":       offset,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func (s *Server) handleUpdateUserByAdmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := mux.Vars(r)["id"]
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			writeJSONError(w, "Invalid user ID", http.StatusBadRequest)
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
			writeJSONError(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert to profile (without password)
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

func (s *Server) handleDeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := mux.Vars(r)["id"]
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			writeJSONError(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Check if trying to delete self
		currentUserID, ok := r.Context().Value(auth.UserIDKey).(*int)
		if ok && *currentUserID == userID {
			writeJSONError(w, "Cannot delete your own account", http.StatusBadRequest)
			return
		}

		if err := database.DeleteUser(r.Context(), s.GetDB(), userID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeJSONError(w, "User not found", http.StatusNotFound)
				return
			}
			writeJSONError(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) handleGetUserStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats, err := database.GetUserStats(r.Context(), s.GetDB())
		if err != nil {
			writeJSONError(w, "Failed to get user stats: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"stats": stats,
		})
	}
}

// Admin-only routes for user management
func (s *Server) RegisterUserRoutes(r *mux.Router) {
	// All user management routes require admin role
	adminRouter := r.PathPrefix("").Subrouter()
	adminRouter.Use(s.jwtManager.AuthMiddleware)
	adminRouter.Use(auth.AdminOnly)

	adminRouter.HandleFunc("/", s.handleListUsers()).Methods("GET")
	adminRouter.HandleFunc("/{id:[0-9]+}", s.handleGetUser()).Methods("GET")
	adminRouter.HandleFunc("/", s.handleCreateUserByAdmin()).Methods("POST")
	adminRouter.HandleFunc("/{id:[0-9]+}", s.handleUpdateUserByAdmin()).Methods("PUT")
	adminRouter.HandleFunc("/{id:[0-9]+}", s.handleDeleteUser()).Methods("DELETE")
	adminRouter.HandleFunc("/stats", s.handleGetUserStats()).Methods("GET")
}