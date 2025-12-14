package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"news-portal-web/api/internal/auth"
	"news-portal-web/api/internal/database"

	"github.com/gorilla/mux"
)

// handleGetArticles returns all articles (with filters)
func (s *Server) handleGetArticles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		filter := database.ArticleFilter{}

		if status := r.URL.Query().Get("status"); status != "" {
			filter.Status = status
		} else {
			// Default to published for public access
			filter.Status = "published"
		}

		if kategoriID := r.URL.Query().Get("kategori_id"); kategoriID != "" {
			if id, err := strconv.Atoi(kategoriID); err == nil {
				filter.KategoriID = id
			}
		}

		// NEW: Support filter by category name
		if kategori := r.URL.Query().Get("kategori"); kategori != "" {
			filter.KategoriName = kategori
		}

		if tagID := r.URL.Query().Get("tag_id"); tagID != "" {
			if id, err := strconv.Atoi(tagID); err == nil {
				filter.TagID = id
			}
		}

		if search := r.URL.Query().Get("search"); search != "" {
			filter.Search = search
		}

		if limit := r.URL.Query().Get("limit"); limit != "" {
			if l, err := strconv.Atoi(limit); err == nil {
				filter.Limit = l
			}
		}

		if offset := r.URL.Query().Get("offset"); offset != "" {
			if o, err := strconv.Atoi(offset); err == nil {
				filter.Offset = o
			}
		}

		articles, err := database.GetAllArticles(s.GetDB(), filter)
		if err != nil {
			writeJSONError(w, "Error fetching articles", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(articles)
	}
}

// handleGetArticleByID returns a single article by ID
func (s *Server) handleGetArticleByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			writeJSONError(w, "Invalid article ID", http.StatusBadRequest)
			return
		}

		article, err := database.GetArticleByID(s.GetDB(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				writeJSONError(w, "Article not found", http.StatusNotFound)
				return
			}
			writeJSONError(w, "Error fetching article", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(article)
	}
}

// handleGetArticleBySlug returns a single article by slug (public)
func (s *Server) handleGetArticleBySlug() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		slug := vars["slug"]

		if slug == "" {
			writeJSONError(w, "Slug is required", http.StatusBadRequest)
			return
		}

		article, err := database.GetPublishedArticleBySlug(s.GetDB(), slug)
		if err != nil {
			if err == sql.ErrNoRows {
				writeJSONError(w, "Article not found", http.StatusNotFound)
				return
			}
			writeJSONError(w, "Error fetching article", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(article)
	}
}

// handleGetArticlesByCategory returns articles by category ID
func (s *Server) handleGetArticlesByCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		kategoriID, err := strconv.Atoi(vars["id"])
		if err != nil {
			writeJSONError(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		limit := 10
		offset := 0

		if l := r.URL.Query().Get("limit"); l != "" {
			if lInt, err := strconv.Atoi(l); err == nil {
				limit = lInt
			}
		}

		if o := r.URL.Query().Get("offset"); o != "" {
			if oInt, err := strconv.Atoi(o); err == nil {
				offset = oInt
			}
		}

		articles, err := database.GetArticlesByCategory(s.GetDB(), kategoriID, limit, offset)
		if err != nil {
			writeJSONError(w, "Error fetching articles", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(articles)
	}
}

// handleCreateArticle creates a new article
func (s *Server) handleCreateArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input database.ArticleInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeJSONError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate input
		if input.Judul == "" {
			writeJSONError(w, "Judul is required", http.StatusBadRequest)
			return
		}

		if input.Konten == "" {
			writeJSONError(w, "Konten is required", http.StatusBadRequest)
			return
		}

		// Get user ID from context
		userID, ok := r.Context().Value(auth.UserIDKey).(int)
		if !ok {
			writeJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		article, err := database.CreateArticle(s.GetDB(), input, userID)
		if err != nil {
			writeJSONError(w, "Error creating article: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(article)
	}
}

// handleUpdateArticle updates an existing article
func (s *Server) handleUpdateArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			writeJSONError(w, "Invalid article ID", http.StatusBadRequest)
			return
		}

		var input database.ArticleInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeJSONError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate input
		if input.Judul == "" {
			writeJSONError(w, "Judul is required", http.StatusBadRequest)
			return
		}

		if input.Konten == "" {
			writeJSONError(w, "Konten is required", http.StatusBadRequest)
			return
		}

		article, err := database.UpdateArticle(s.GetDB(), id, input)
		if err != nil {
			if err == sql.ErrNoRows {
				writeJSONError(w, "Article not found", http.StatusNotFound)
				return
			}
			writeJSONError(w, "Error updating article: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(article)
	}
}

// handleDeleteArticle deletes an article
func (s *Server) handleDeleteArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			writeJSONError(w, "Invalid article ID", http.StatusBadRequest)
			return
		}

		err = database.DeleteArticle(s.GetDB(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				writeJSONError(w, "Article not found", http.StatusNotFound)
				return
			}
			writeJSONError(w, "Error deleting article", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Article deleted successfully"})
	}
}

// ========================================
// ROUTE REGISTRATION
// ========================================

// RegisterPublicArticleRoutes registers public article routes
func (s *Server) RegisterPublicArticleRoutes(r *mux.Router) {
	r.HandleFunc("/articles", s.handleGetArticles()).Methods("GET")
	r.HandleFunc("/articles/{id:[0-9]+}", s.handleGetArticleByID()).Methods("GET")
	r.HandleFunc("/articles/slug/{slug}", s.handleGetArticleBySlug()).Methods("GET")
}

// RegisterEditorArticleRoutes registers editor article routes
func (s *Server) RegisterEditorArticleRoutes(r *mux.Router) {
	r.HandleFunc("/articles", s.handleCreateArticle()).Methods("POST")
	r.HandleFunc("/articles/{id:[0-9]+}", s.handleUpdateArticle()).Methods("PUT")
	r.HandleFunc("/articles/{id:[0-9]+}", s.handleDeleteArticle()).Methods("DELETE")
}
