package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"news-portal-web/api/internal/database"

	"github.com/gorilla/mux"
)

func (s *Server) handleCreateCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req database.CategoryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := validateCategoryRequest(&req); err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if category already exists
		exists, err := database.IsCategoryExists(r.Context(), s.GetDB(), req.NamaKategori)
		if err != nil {
			writeJSONError(w, "Failed to check category existence: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if exists {
			writeJSONError(w, "Category already exists", http.StatusConflict)
			return
		}

		category, err := database.CreateCategory(r.Context(), s.GetDB(), &req)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				writeJSONError(w, "Category already exists", http.StatusConflict)
				return
			}
			writeJSONError(w, "Failed to create category: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(category)
	}
}

func (s *Server) handleGetCategoryByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categoryID, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			writeJSONError(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		category, err := database.GetCategoryByID(r.Context(), s.GetDB(), categoryID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found") {
				writeJSONError(w, "Category not found", http.StatusNotFound)
			} else {
				writeJSONError(w, "Failed to get category: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(category)
	}
}

func (s *Server) handleGetCategories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if requesting categories with article count
		withCount := r.URL.Query().Get("with_count")

		if withCount == "true" {
			categories, err := database.ListCategoriesWithArticleCount(r.Context(), s.GetDB())
			if err != nil {
				writeJSONError(w, "Failed to fetch categories: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"categories": categories,
			})
			return
		}

		categories, err := database.ListCategories(r.Context(), s.GetDB())
		if err != nil {
			writeJSONError(w, "Failed to fetch categories: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"categories": categories,
		})
	}
}

func (s *Server) handleUpdateCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categoryIDStr := mux.Vars(r)["id"]
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			writeJSONError(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		var req database.CategoryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := validateCategoryRequest(&req); err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if new name already exists (excluding current category)
		existing, err := database.GetCategoryByName(r.Context(), s.GetDB(), req.NamaKategori)
		if err == nil && existing.KategoriID != categoryID {
			writeJSONError(w, "Category name already exists", http.StatusConflict)
			return
		}

		category, err := database.UpdateCategory(r.Context(), s.GetDB(), categoryID, &req)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				writeJSONError(w, "Category not found", http.StatusNotFound)
				return
			}
			if strings.Contains(err.Error(), "duplicate") {
				writeJSONError(w, "Category name already exists", http.StatusConflict)
				return
			}
			writeJSONError(w, "Failed to update category: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(category)
	}
}

func (s *Server) handleDeleteCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categoryIDStr := mux.Vars(r)["id"]
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			writeJSONError(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		// Check for force delete parameter
		force := r.URL.Query().Get("force")

		if force == "true" {
			// Force delete - removes category and all its relations
			err = database.ForceDeleteCategory(r.Context(), s.GetDB(), categoryID)
		} else {
			// Safe delete - only delete if no articles are using this category
			err = database.DeleteCategory(r.Context(), s.GetDB(), categoryID)
		}

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeJSONError(w, "Category not found", http.StatusNotFound)
				return
			}
			if strings.Contains(err.Error(), "cannot delete category that has articles") {
				writeJSONError(w, "Cannot delete category that has articles. Use force=true to delete anyway.", http.StatusConflict)
				return
			}
			writeJSONError(w, "Failed to delete category: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) RegisterCategoryRoutes(r *mux.Router) {
	// Public routes
	r.HandleFunc("/", s.handleGetCategories()).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}", s.handleGetCategoryByID()).Methods("GET")

	// Protected routes (requires authentication)
	r.HandleFunc("/", s.handleCreateCategory()).Methods("POST")
	r.HandleFunc("/{id:[0-9]+}", s.handleUpdateCategory()).Methods("PUT")
	r.HandleFunc("/{id:[0-9]+}", s.handleDeleteCategory()).Methods("DELETE")
}

func validateCategoryRequest(req *database.CategoryRequest) error {
	if len(strings.TrimSpace(req.NamaKategori)) == 0 {
		return errors.New("category name is required")
	}

	if len(req.NamaKategori) < 2 || len(req.NamaKategori) > 100 {
		return errors.New("category name must be between 2 and 100 characters")
	}

	// Trim spaces and normalize
	req.NamaKategori = strings.TrimSpace(req.NamaKategori)

	return nil
}

// ========================================
// ROUTE REGISTRATION
// ========================================

// RegisterPublicCategoryRoutes registers public category routes
func (s *Server) RegisterPublicCategoryRoutes(r *mux.Router) {
	r.HandleFunc("/categories", s.handleGetCategories()).Methods("GET")
	r.HandleFunc("/categories/{id:[0-9]+}", s.handleGetCategoryByID()).Methods("GET")
}

// RegisterAdminCategoryRoutes registers admin category routes
func (s *Server) RegisterAdminCategoryRoutes(r *mux.Router) {
	r.HandleFunc("/categories", s.handleCreateCategory()).Methods("POST")
	r.HandleFunc("/categories/{id:[0-9]+}", s.handleUpdateCategory()).Methods("PUT")
	r.HandleFunc("/categories/{id:[0-9]+}", s.handleDeleteCategory()).Methods("DELETE")
}
