package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"news-portal/api/internal/database"

	"github.com/gorilla/mux"
)

func (s *Server) handleCreateComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context (optional for comments)
		userID, _ := r.Context().Value("user_id").(*int)

		var req database.CommentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := validateCommentRequest(&req); err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		comment, err := database.CreateComment(r.Context(), s.GetDB(), userID, &req)
		if err != nil {
			if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "not published") {
				writeJSONError(w, "Article not found or not published", http.StatusNotFound)
				return
			}
			writeJSONError(w, "Failed to create comment: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(comment)
	}
}

func (s *Server) handleGetComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		commentIDStr := vars["id"]

		if commentIDStr != "" {
			commentID, err := strconv.Atoi(commentIDStr)
			if err != nil {
				writeJSONError(w, "Invalid comment ID", http.StatusBadRequest)
				return
			}

			comment, err := database.GetCommentByID(r.Context(), s.GetDB(), commentID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found") {
					writeJSONError(w, "Comment not found", http.StatusNotFound)
				} else {
					writeJSONError(w, "Failed to get comment: "+err.Error(), http.StatusInternalServerError)
				}
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(comment)
			return
		}

		// List comments
		s.handleListComments()(w, r)
	}
}

func (s *Server) handleListComments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		articleIDStr := r.URL.Query().Get("article_id")
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		// If article_id is provided, get comments for specific article
		if articleIDStr != "" {
			articleID, err := strconv.Atoi(articleIDStr)
			if err != nil {
				writeJSONError(w, "Invalid article ID", http.StatusBadRequest)
				return
			}

			comments, err := database.ListCommentsByArticle(r.Context(), s.GetDB(), articleID)
			if err != nil {
				writeJSONError(w, "Failed to fetch comments: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"comments": comments,
				"count":    len(comments),
			})
			return
		}

		// List all comments with pagination (admin functionality)
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

		comments, totalCount, err := database.ListAllComments(r.Context(), s.GetDB(), limit, offset)
		if err != nil {
			writeJSONError(w, "Failed to fetch comments: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Calculate pagination info
		totalPages := (totalCount + limit - 1) / limit
		currentPage := (offset / limit) + 1

		response := map[string]interface{}{
			"comments": comments,
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

func (s *Server) handleUpdateComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("user_id").(*int)

		commentIDStr := mux.Vars(r)["id"]
		commentID, err := strconv.Atoi(commentIDStr)
		if err != nil {
			writeJSONError(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		var req struct {
			Konten string `json:"konten"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := validateCommentContent(req.Konten); err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		comment, err := database.UpdateComment(r.Context(), s.GetDB(), commentID, userID, req.Konten)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				writeJSONError(w, "Comment not found", http.StatusNotFound)
				return
			}
			if strings.Contains(err.Error(), "unauthorized") {
				writeJSONError(w, "You can only edit your own comments", http.StatusForbidden)
				return
			}
			writeJSONError(w, "Failed to update comment: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(comment)
	}
}

func (s *Server) handleDeleteComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("user_id").(*int)

		commentIDStr := mux.Vars(r)["id"]
		commentID, err := strconv.Atoi(commentIDStr)
		if err != nil {
			writeJSONError(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		if err := database.DeleteComment(r.Context(), s.GetDB(), commentID, userID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeJSONError(w, "Comment not found", http.StatusNotFound)
				return
			}
			if strings.Contains(err.Error(), "unauthorized") {
				writeJSONError(w, "You can only delete your own comments", http.StatusForbidden)
				return
			}
			writeJSONError(w, "Failed to delete comment: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) handleGetCommentStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		totalCount, err := database.GetTotalCommentCount(r.Context(), s.GetDB())
		if err != nil {
			writeJSONError(w, "Failed to get comment stats: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"total_comments": totalCount,
		})
	}
}

func (s *Server) RegisterCommentRoutes(r *mux.Router) {
	// Public routes
	r.HandleFunc("/", s.handleGetComment()).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}", s.handleGetComment()).Methods("GET")
	r.HandleFunc("/", s.handleCreateComment()).Methods("POST")

	// User routes (requires authentication but optional)
	r.HandleFunc("/{id:[0-9]+}", s.handleUpdateComment()).Methods("PUT")
	r.HandleFunc("/{id:[0-9]+}", s.handleDeleteComment()).Methods("DELETE")

	// Stats route
	r.HandleFunc("/stats", s.handleGetCommentStats()).Methods("GET")
}

func validateCommentRequest(req *database.CommentRequest) error {
	if req.ArtikelID <= 0 {
		return errors.New("valid article ID is required")
	}

	return validateCommentContent(req.Konten)
}

func validateCommentContent(konten string) error {
	if len(strings.TrimSpace(konten)) == 0 {
		return errors.New("comment content is required")
	}

	if len(konten) < 1 || len(konten) > 2000 {
		return errors.New("comment content must be between 1 and 2000 characters")
	}

	return nil
}
