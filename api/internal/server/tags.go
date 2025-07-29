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

func (s *Server) handleCreateTag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req database.TagRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := validateTagRequest(&req); err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if tag already exists
		exists, err := database.IsTagExists(r.Context(), s.GetDB(), req.NamaTag)
		if err != nil {
			writeJSONError(w, "Failed to check tag existence: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if exists {
			writeJSONError(w, "Tag already exists", http.StatusConflict)
			return
		}

		tag, err := database.CreateTag(r.Context(), s.GetDB(), &req)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				writeJSONError(w, "Tag already exists", http.StatusConflict)
				return
			}
			writeJSONError(w, "Failed to create tag: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tag)
	}
}

func (s *Server) handleGetTag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tagIDStr := vars["id"]

		if tagIDStr != "" {
			tagID, err := strconv.Atoi(tagIDStr)
			if err != nil {
				writeJSONError(w, "Invalid tag ID", http.StatusBadRequest)
				return
			}

			tag, err := database.GetTagByID(r.Context(), s.GetDB(), tagID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found") {
					writeJSONError(w, "Tag not found", http.StatusNotFound)
				} else {
					writeJSONError(w, "Failed to get tag: "+err.Error(), http.StatusInternalServerError)
				}
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tag)
			return
		}

		// List all tags
		s.handleListTags()(w, r)
	}
}

func (s *Server) handleListTags() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		withCount := r.URL.Query().Get("with_count")
		popular := r.URL.Query().Get("popular")
		search := r.URL.Query().Get("search")
		limitStr := r.URL.Query().Get("limit")

		// Search tags
		if search != "" {
			tags, err := database.SearchTags(r.Context(), s.GetDB(), search)
			if err != nil {
				writeJSONError(w, "Failed to search tags: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"tags": tags,
			})
			return
		}

		// Popular tags
		if popular == "true" {
			limit := 10
			if limitStr != "" {
				if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
					limit = l
				}
			}

			tags, err := database.ListPopularTags(r.Context(), s.GetDB(), limit)
			if err != nil {
				writeJSONError(w, "Failed to fetch popular tags: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"tags": tags,
			})
			return
		}

		// Tags with article count
		if withCount == "true" {
			tags, err := database.ListTagsWithArticleCount(r.Context(), s.GetDB())
			if err != nil {
				writeJSONError(w, "Failed to fetch tags: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"tags": tags,
			})
			return
		}

		// Default: list all tags
		tags, err := database.ListTags(r.Context(), s.GetDB())
		if err != nil {
			writeJSONError(w, "Failed to fetch tags: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"tags": tags,
		})
	}
}

func (s *Server) handleUpdateTag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tagIDStr := mux.Vars(r)["id"]
		tagID, err := strconv.Atoi(tagIDStr)
		if err != nil {
			writeJSONError(w, "Invalid tag ID", http.StatusBadRequest)
			return
		}

		var req database.TagRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := validateTagRequest(&req); err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if new name already exists (excluding current tag)
		existing, err := database.GetTagByName(r.Context(), s.GetDB(), req.NamaTag)
		if err == nil && existing.TagID != tagID {
			writeJSONError(w, "Tag name already exists", http.StatusConflict)
			return
		}

		tag, err := database.UpdateTag(r.Context(), s.GetDB(), tagID, &req)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				writeJSONError(w, "Tag not found", http.StatusNotFound)
				return
			}
			if strings.Contains(err.Error(), "duplicate") {
				writeJSONError(w, "Tag name already exists", http.StatusConflict)
				return
			}
			writeJSONError(w, "Failed to update tag: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tag)
	}
}

func (s *Server) handleDeleteTag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tagIDStr := mux.Vars(r)["id"]
		tagID, err := strconv.Atoi(tagIDStr)
		if err != nil {
			writeJSONError(w, "Invalid tag ID", http.StatusBadRequest)
			return
		}

		// Check for force delete parameter
		force := r.URL.Query().Get("force")

		if force == "true" {
			// Force delete - removes tag and all its relations
			err = database.ForceDeleteTag(r.Context(), s.GetDB(), tagID)
		} else {
			// Safe delete - only delete if no articles are using this tag
			err = database.DeleteTag(r.Context(), s.GetDB(), tagID)
		}

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeJSONError(w, "Tag not found", http.StatusNotFound)
				return
			}
			if strings.Contains(err.Error(), "cannot delete tag that has articles") {
				writeJSONError(w, "Cannot delete tag that has articles. Use force=true to delete anyway.", http.StatusConflict)
				return
			}
			writeJSONError(w, "Failed to delete tag: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) handleCreateMultipleTags() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			TagNames []string `json:"tag_names"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		if len(req.TagNames) == 0 {
			writeJSONError(w, "Tag names are required", http.StatusBadRequest)
			return
		}

		// Validate each tag name
		for _, tagName := range req.TagNames {
			tagReq := database.TagRequest{NamaTag: tagName}
			if err := validateTagRequest(&tagReq); err != nil {
				writeJSONError(w, "Invalid tag '"+tagName+"': "+err.Error(), http.StatusBadRequest)
				return
			}
		}

		tagIDs, err := database.GetOrCreateTags(r.Context(), s.GetDB(), req.TagNames)
		if err != nil {
			writeJSONError(w, "Failed to create tags: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"tag_ids": tagIDs,
			"message": "Tags created successfully",
		})
	}
}

func (s *Server) RegisterTagRoutes(r *mux.Router) {
	// Public routes
	r.HandleFunc("/", s.handleGetTag()).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}", s.handleGetTag()).Methods("GET")

	// Protected routes (requires authentication)
	r.HandleFunc("/", s.handleCreateTag()).Methods("POST")
	r.HandleFunc("/bulk", s.handleCreateMultipleTags()).Methods("POST")
	r.HandleFunc("/{id:[0-9]+}", s.handleUpdateTag()).Methods("PUT")
	r.HandleFunc("/{id:[0-9]+}", s.handleDeleteTag()).Methods("DELETE")
}

func validateTagRequest(req *database.TagRequest) error {
	if len(strings.TrimSpace(req.NamaTag)) == 0 {
		return errors.New("tag name is required")
	}

	if len(req.NamaTag) < 2 || len(req.NamaTag) > 50 {
		return errors.New("tag name must be between 2 and 50 characters")
	}

	// Check for invalid characters (only allow letters, numbers, spaces, hyphens)
	for _, char := range req.NamaTag {
		if !(char >= 'a' && char <= 'z') &&
			!(char >= 'A' && char <= 'Z') &&
			!(char >= '0' && char <= '9') &&
			char != ' ' && char != '-' && char != '_' {
			return errors.New("tag name can only contain letters, numbers, spaces, hyphens, and underscores")
		}
	}

	// Trim spaces and normalize
	req.NamaTag = strings.TrimSpace(req.NamaTag)

	return nil
}
