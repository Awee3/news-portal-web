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

func (s *Server) handleCreateArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context (set by auth middleware)
		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req database.ArticleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Basic validation
		if err := validateArticleRequest(&req); err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Use s.GetDB() instead of s.db
		article, err := database.CreateArticle(r.Context(), s.GetDB(), userID, &req)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				writeJSONError(w, "Article with this title already exists", http.StatusConflict)
				return
			}
			writeJSONError(w, "Failed to create article: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(article)
	}
}

func (s *Server) handleGetArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// Check if it's getting by slug or ID
		if slug := vars["slug"]; slug != "" {
			article, err := database.GetArticleBySlug(r.Context(), s.GetDB(), slug)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found") {
					writeJSONError(w, "Article not found", http.StatusNotFound)
				} else {
					writeJSONError(w, "Failed to get article: "+err.Error(), http.StatusInternalServerError)
				}
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(article)
			return
		}

		if articleIDStr := vars["id"]; articleIDStr != "" {
			articleID, err := strconv.Atoi(articleIDStr)
			if err != nil {
				writeJSONError(w, "Invalid article ID", http.StatusBadRequest)
				return
			}

			article, err := database.GetArticleByID(r.Context(), s.GetDB(), articleID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found") {
					writeJSONError(w, "Article not found", http.StatusNotFound)
				} else {
					writeJSONError(w, "Failed to get article: "+err.Error(), http.StatusInternalServerError)
				}
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(article)
			return
		}

		// List articles with pagination and filters
		s.handleListArticles()(w, r)
	}
}

func (s *Server) handleListArticles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")
		categoryIDStr := r.URL.Query().Get("category_id")
		status := r.URL.Query().Get("status")

		// Default values
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

		var categoryID *int
		if categoryIDStr != "" {
			if cid, err := strconv.Atoi(categoryIDStr); err == nil {
				categoryID = &cid
			}
		}

		// Default to published for public access
		if status == "" {
			status = "published"
		}

		articles, totalCount, err := database.ListArticles(r.Context(), s.GetDB(), limit, offset, categoryID, status)
		if err != nil {
			writeJSONError(w, "Failed to fetch articles: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Calculate pagination info
		totalPages := (totalCount + limit - 1) / limit
		currentPage := (offset / limit) + 1

		response := map[string]interface{}{
			"articles": articles,
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

func (s *Server) handleUpdateArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		articleIDStr := mux.Vars(r)["id"]
		articleID, err := strconv.Atoi(articleIDStr)
		if err != nil {
			writeJSONError(w, "Invalid article ID", http.StatusBadRequest)
			return
		}

		var req database.ArticleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := validateArticleRequest(&req); err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		article, err := database.UpdateArticle(r.Context(), s.GetDB(), articleID, userID, &req)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				writeJSONError(w, "Article not found", http.StatusNotFound)
				return
			}
			if strings.Contains(err.Error(), "unauthorized") {
				writeJSONError(w, "You can only edit your own articles", http.StatusForbidden)
				return
			}
			writeJSONError(w, "Failed to update article: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(article)
	}
}

func (s *Server) handleDeleteArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		articleIDStr := mux.Vars(r)["id"]
		articleID, err := strconv.Atoi(articleIDStr)
		if err != nil {
			writeJSONError(w, "Invalid article ID", http.StatusBadRequest)
			return
		}

		if err := database.DeleteArticle(r.Context(), s.GetDB(), articleID, userID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeJSONError(w, "Article not found", http.StatusNotFound)
				return
			}
			if strings.Contains(err.Error(), "unauthorized") {
				writeJSONError(w, "You can only delete your own articles", http.StatusForbidden)
				return
			}
			writeJSONError(w, "Failed to delete article: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) RegisterArticleRoutes(r *mux.Router) {
	// Public routes
	r.HandleFunc("/", s.handleGetArticle()).Methods("GET")
	r.HandleFunc("/slug/{slug}", s.handleGetArticle()).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}", s.handleGetArticle()).Methods("GET")

	// Protected routes (requires authentication)
	r.HandleFunc("/", s.handleCreateArticle()).Methods("POST")
	r.HandleFunc("/{id:[0-9]+}", s.handleUpdateArticle()).Methods("PUT")
	r.HandleFunc("/{id:[0-9]+}", s.handleDeleteArticle()).Methods("DELETE")
}

func validateArticleRequest(req *database.ArticleRequest) error {
	if len(req.Judul) < 10 || len(req.Judul) > 200 {
		return errors.New("title must be between 10 and 200 characters")
	}

	if len(req.Konten) < 50 {
		return errors.New("content must be at least 50 characters")
	}

	if req.Status != "draft" && req.Status != "published" {
		return errors.New("status must be either 'draft' or 'published'")
	}

	if req.MetaTitle != nil && len(*req.MetaTitle) > 60 {
		return errors.New("meta title must not exceed 60 characters")
	}

	if req.MetaDescription != nil && len(*req.MetaDescription) > 160 {
		return errors.New("meta description must not exceed 160 characters")
	}

	return nil
}

