package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"news-portal-web/api/internal/database"

	"github.com/gorilla/mux"
)

// ========================================
// AUTHENTICATED USER HANDLERS
// ========================================

// handleGetCurrentUser - GET /api/v1/users/me
func (s *Server) handleGetCurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey).(*Claims)
		if !ok || claims == nil {
			writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := database.GetUserByID(s.GetDB(), claims.UserID)
		if err != nil {
			writeJSONError(w, "User tidak ditemukan", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user.ToPublic())
	}
}

// UpdateUserRequest - Request body untuk update user
type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// handleUpdateCurrentUser - PUT /api/v1/users/me
func (s *Server) handleUpdateCurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey).(*Claims)
		if !ok || claims == nil {
			writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validasi
		if req.Username == "" {
			writeJSONError(w, "Username harus diisi", http.StatusBadRequest)
			return
		}

		if req.Email == "" {
			writeJSONError(w, "Email harus diisi", http.StatusBadRequest)
			return
		}

		if !isValidEmail(req.Email) {
			writeJSONError(w, "Format email tidak valid", http.StatusBadRequest)
			return
		}

		// Check username uniqueness
		exists, err := database.CheckUsernameExists(s.GetDB(), req.Username, claims.UserID)
		if err != nil {
			writeJSONError(w, "Error checking username", http.StatusInternalServerError)
			return
		}
		if exists {
			writeJSONError(w, "Username sudah digunakan", http.StatusConflict)
			return
		}

		// Check email uniqueness
		exists, err = database.CheckEmailExists(s.GetDB(), req.Email, claims.UserID)
		if err != nil {
			writeJSONError(w, "Error checking email", http.StatusInternalServerError)
			return
		}
		if exists {
			writeJSONError(w, "Email sudah digunakan", http.StatusConflict)
			return
		}

		// Gunakan UpdateUserBasic bukan UpdateUser
		user, err := database.UpdateUserBasic(s.GetDB(), claims.UserID, req.Username, req.Email)
		if err != nil {
			writeJSONError(w, "Error updating user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user.ToPublic())
	}
}

// ChangePasswordRequest - Request body untuk ganti password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// handleChangePassword - PUT /api/v1/users/me/password
func (s *Server) handleChangePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey).(*Claims)
		if !ok || claims == nil {
			writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req ChangePasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.CurrentPassword == "" || req.NewPassword == "" {
			writeJSONError(w, "Password lama dan baru harus diisi", http.StatusBadRequest)
			return
		}

		if len(req.NewPassword) < 8 {
			writeJSONError(w, "Password baru minimal 8 karakter", http.StatusBadRequest)
			return
		}

		// Get current user
		user, err := database.GetUserByID(s.GetDB(), claims.UserID)
		if err != nil {
			writeJSONError(w, "User tidak ditemukan", http.StatusNotFound)
			return
		}

		// Verify current password
		if !database.VerifyPassword(user.Password, req.CurrentPassword) {
			writeJSONError(w, "Password lama tidak sesuai", http.StatusBadRequest)
			return
		}

		// Update password
		err = database.UpdateUserPassword(s.GetDB(), claims.UserID, req.NewPassword)
		if err != nil {
			writeJSONError(w, "Error updating password", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Password berhasil diperbarui",
		})
	}
}

// ========================================
// ADMIN HANDLERS
// ========================================

// handleGetAllUsers - GET /api/v1/admin/users
func (s *Server) handleGetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := database.GetAllUsers(s.GetDB())
		if err != nil {
			writeJSONError(w, "Error fetching users", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

// handleGetUserByID - GET /api/v1/admin/users/{id}
func (s *Server) handleGetUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, err := strconv.Atoi(vars["id"])
		if err != nil {
			writeJSONError(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		user, err := database.GetUserByID(s.GetDB(), userID)
		if err != nil {
			writeJSONError(w, "User tidak ditemukan", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user.ToPublic())
	}
}

// UpdateUserRoleRequest - Request body untuk update role
type UpdateUserRoleRequest struct {
	Role string `json:"role"`
}

// handleUpdateUserRole - PUT /api/v1/admin/users/{id}/role
func (s *Server) handleUpdateUserRole() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, err := strconv.Atoi(vars["id"])
		if err != nil {
			writeJSONError(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Check if user exists
		_, err = database.GetUserByID(s.GetDB(), userID)
		if err != nil {
			writeJSONError(w, "User tidak ditemukan", http.StatusNotFound)
			return
		}

		var req UpdateUserRoleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate role
		if req.Role != "admin" && req.Role != "editor" && req.Role != "user" {
			writeJSONError(w, "Role tidak valid (admin, editor, user)", http.StatusBadRequest)
			return
		}

		err = database.UpdateUserRole(s.GetDB(), userID, req.Role)
		if err != nil {
			writeJSONError(w, "Error updating role", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Role berhasil diperbarui",
			"role":    req.Role,
		})
	}
}

// handleDeleteUser - DELETE /api/v1/admin/users/{id}
func (s *Server) handleDeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey).(*Claims)
		if !ok || claims == nil {
			writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		userID, err := strconv.Atoi(vars["id"])
		if err != nil {
			writeJSONError(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Prevent self-deletion
		if userID == claims.UserID {
			writeJSONError(w, "Tidak bisa menghapus akun sendiri", http.StatusBadRequest)
			return
		}

		// Check if user exists
		_, err = database.GetUserByID(s.GetDB(), userID)
		if err != nil {
			writeJSONError(w, "User tidak ditemukan", http.StatusNotFound)
			return
		}

		err = database.DeleteUser(s.GetDB(), userID)
		if err != nil {
			writeJSONError(w, "Error deleting user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User berhasil dihapus",
		})
	}
}

// ========================================
// ROUTE REGISTRATION
// ========================================

// RegisterAdminUserRoutes - Register admin user routes
func (s *Server) RegisterAdminUserRoutes(r *mux.Router) {
	r.HandleFunc("/users", s.handleGetAllUsers()).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", s.handleGetUserByID()).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}/role", s.handleUpdateUserRole()).Methods("PUT")
	r.HandleFunc("/users/{id:[0-9]+}", s.handleDeleteUser()).Methods("DELETE")
}
