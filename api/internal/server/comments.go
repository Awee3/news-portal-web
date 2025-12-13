package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"news-portal-web/api/internal/database"

	"github.com/gorilla/mux"
)

// ========================================
// PUBLIC HANDLERS
// ========================================

// handleGetArticleComments - GET /api/v1/articles/{id}/comments
// Mendapatkan semua komentar yang sudah approved untuk artikel tertentu
func (s *Server) handleGetArticleComments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		articleID, err := strconv.Atoi(vars["id"])
		if err != nil {
			writeJSONError(w, "Invalid article ID", http.StatusBadRequest)
			return
		}

		// Hanya tampilkan komentar yang sudah approved untuk public
		comments, err := database.GetCommentsByArticleID(s.GetDB(), articleID, "approved")
		if err != nil {
			writeJSONError(w, "Error fetching comments", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	}
}

// CreateCommentRequest - Request body untuk membuat komentar
type CreateCommentRequest struct {
	Konten       string `json:"konten"`
	NamaPengguna string `json:"nama_pengguna,omitempty"` // untuk anonymous
}

// handleCreateComment - POST /api/v1/articles/{id}/comments
// Membuat komentar baru (bisa anonymous atau authenticated)
func (s *Server) handleCreateComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		articleID, err := strconv.Atoi(vars["id"])
		if err != nil {
			writeJSONError(w, "Invalid article ID", http.StatusBadRequest)
			return
		}

		var req CreateCommentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validasi konten
		if req.Konten == "" {
			writeJSONError(w, "Konten komentar harus diisi", http.StatusBadRequest)
			return
		}

		if len(req.Konten) > 2000 {
			writeJSONError(w, "Konten komentar maksimal 2000 karakter", http.StatusBadRequest)
			return
		}

		// Check apakah artikel exists
		_, err = database.GetArticleByID(s.GetDB(), articleID)
		if err != nil {
			writeJSONError(w, "Artikel tidak ditemukan", http.StatusNotFound)
			return
		}

		// Cek apakah user authenticated (optional auth)
		var userID *int
		var namaPengguna *string

		if claims, ok := r.Context().Value(ClaimsKey).(*Claims); ok && claims != nil {
			userID = &claims.UserID
		} else {
			// Anonymous comment - nama_pengguna required
			if req.NamaPengguna == "" {
				req.NamaPengguna = "Anonymous"
			}
			namaPengguna = &req.NamaPengguna
		}

		comment := &database.Comment{
			Konten:       req.Konten,
			NamaPengguna: namaPengguna,
			Status:       "pending", // Default pending, perlu moderasi
			UserID:       userID,
			ArtikelID:    articleID,
		}

		createdComment, err := database.CreateCommentSimple(s.GetDB(), comment)
		if err != nil {
			writeJSONError(w, "Error creating comment", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdComment)
	}
}

// ========================================
// AUTHENTICATED USER HANDLERS
// ========================================

// handleGetUserComments - GET /api/v1/users/me/comments
// Mendapatkan komentar milik user yang sedang login
func (s *Server) handleGetUserComments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey).(*Claims)
		if !ok || claims == nil {
			writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		comments, err := database.GetCommentsByUserID(s.GetDB(), claims.UserID)
		if err != nil {
			writeJSONError(w, "Error fetching comments", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	}
}

// handleUpdateUserComment - PUT /api/v1/users/me/comments/{id}
// User mengupdate komentar miliknya sendiri
func (s *Server) handleUpdateUserComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey).(*Claims)
		if !ok || claims == nil {
			writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		commentID, err := strconv.Atoi(vars["id"])
		if err != nil {
			writeJSONError(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		// Cek ownership
		existingComment, err := database.GetCommentByIDSimple(s.GetDB(), commentID)
		if err != nil {
			writeJSONError(w, "Komentar tidak ditemukan", http.StatusNotFound)
			return
		}

		if existingComment.UserID == nil || *existingComment.UserID != claims.UserID {
			writeJSONError(w, "Anda tidak memiliki akses untuk mengubah komentar ini", http.StatusForbidden)
			return
		}

		var req CreateCommentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Konten == "" {
			writeJSONError(w, "Konten komentar harus diisi", http.StatusBadRequest)
			return
		}

		// Update komentar - status kembali ke pending untuk re-moderasi
		updatedComment, err := database.UpdateCommentSimple(s.GetDB(), commentID, req.Konten, "pending")
		if err != nil {
			writeJSONError(w, "Error updating comment", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedComment)
	}
}

// handleDeleteUserComment - DELETE /api/v1/users/me/comments/{id}
// User menghapus komentar miliknya sendiri
func (s *Server) handleDeleteUserComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey).(*Claims)
		if !ok || claims == nil {
			writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		commentID, err := strconv.Atoi(vars["id"])
		if err != nil {
			writeJSONError(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		// Cek ownership
		existingComment, err := database.GetCommentByIDSimple(s.GetDB(), commentID)
		if err != nil {
			writeJSONError(w, "Komentar tidak ditemukan", http.StatusNotFound)
			return
		}

		if existingComment.UserID == nil || *existingComment.UserID != claims.UserID {
			writeJSONError(w, "Anda tidak memiliki akses untuk menghapus komentar ini", http.StatusForbidden)
			return
		}

		err = database.DeleteCommentSimple(s.GetDB(), commentID)
		if err != nil {
			writeJSONError(w, "Error deleting comment", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Komentar berhasil dihapus",
		})
	}
}

// ========================================
// ADMIN/MODERATOR HANDLERS
// ========================================

// handleGetAllComments - GET /api/v1/admin/comments
// Admin mendapatkan semua komentar dengan filter status
func (s *Server) handleGetAllComments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status") // pending, approved, rejected, atau kosong untuk semua

		// Parse pagination
		limit := 50 // default
		offset := 0
		if l := r.URL.Query().Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil {
				limit = parsed
			}
		}
		if o := r.URL.Query().Get("offset"); o != "" {
			if parsed, err := strconv.Atoi(o); err == nil {
				offset = parsed
			}
		}

		comments, err := database.GetAllComments(s.GetDB(), status, limit, offset)
		if err != nil {
			writeJSONError(w, "Error fetching comments", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	}
}

// ModerateCommentRequest - Request body untuk moderasi komentar
type ModerateCommentRequest struct {
	Status string `json:"status"` // approved atau rejected
}

// handleModerateComment - PUT /api/v1/admin/comments/{id}/moderate
// Admin menyetujui atau menolak komentar
func (s *Server) handleModerateComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		commentID, err := strconv.Atoi(vars["id"])
		if err != nil {
			writeJSONError(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		var req ModerateCommentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validasi status
		if req.Status != "approved" && req.Status != "rejected" {
			writeJSONError(w, "Status harus 'approved' atau 'rejected'", http.StatusBadRequest)
			return
		}

		// Cek apakah komentar exists
		_, err = database.GetCommentByIDSimple(s.GetDB(), commentID)
		if err != nil {
			writeJSONError(w, "Komentar tidak ditemukan", http.StatusNotFound)
			return
		}

		// Update status
		updatedComment, err := database.UpdateCommentStatus(s.GetDB(), commentID, req.Status)
		if err != nil {
			writeJSONError(w, "Error moderating comment", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedComment)
	}
}

// handleAdminDeleteComment - DELETE /api/v1/admin/comments/{id}
// Admin menghapus komentar apapun
func (s *Server) handleAdminDeleteComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		commentID, err := strconv.Atoi(vars["id"])
		if err != nil {
			writeJSONError(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		// Cek apakah komentar exists
		_, err = database.GetCommentByIDSimple(s.GetDB(), commentID)
		if err != nil {
			writeJSONError(w, "Komentar tidak ditemukan", http.StatusNotFound)
			return
		}

		err = database.DeleteCommentSimple(s.GetDB(), commentID)
		if err != nil {
			writeJSONError(w, "Error deleting comment", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Komentar berhasil dihapus",
		})
	}
}

// ========================================
// ROUTE REGISTRATION
// ========================================

// RegisterPublicCommentRoutes - Register public comment routes
func (s *Server) RegisterPublicCommentRoutes(r *mux.Router) {
	r.HandleFunc("/articles/{id:[0-9]+}/comments", s.handleGetArticleComments()).Methods("GET")
	r.HandleFunc("/articles/{id:[0-9]+}/comments", s.handleCreateComment()).Methods("POST")
}

// RegisterUserCommentRoutes - Register authenticated user comment routes
func (s *Server) RegisterUserCommentRoutes(r *mux.Router) {
	r.HandleFunc("/users/me/comments", s.handleGetUserComments()).Methods("GET")
	r.HandleFunc("/users/me/comments/{id:[0-9]+}", s.handleUpdateUserComment()).Methods("PUT")
	r.HandleFunc("/users/me/comments/{id:[0-9]+}", s.handleDeleteUserComment()).Methods("DELETE")
}

// RegisterAdminCommentRoutes - Register admin comment routes
func (s *Server) RegisterAdminCommentRoutes(r *mux.Router) {
	r.HandleFunc("/comments", s.handleGetAllComments()).Methods("GET")
	r.HandleFunc("/comments/{id:[0-9]+}/moderate", s.handleModerateComment()).Methods("PUT")
	r.HandleFunc("/comments/{id:[0-9]+}", s.handleAdminDeleteComment()).Methods("DELETE")
}
