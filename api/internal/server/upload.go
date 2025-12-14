package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	// "strings"
	"time"
)

// handleUpload - POST /api/v1/editor/upload
func (s *Server) handleUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Max 10MB
		r.ParseMultipartForm(10 << 20)

		file, header, err := r.FormFile("file")
		if err != nil {
			writeJSONError(w, "File tidak ditemukan", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Validate file type
		contentType := header.Header.Get("Content-Type")
		allowedTypes := map[string]string{
			"image/jpeg": ".jpg",
			"image/png":  ".png",
			"image/gif":  ".gif",
			"image/webp": ".webp",
		}

		ext, ok := allowedTypes[contentType]
		if !ok {
			writeJSONError(w, "Format file tidak didukung", http.StatusBadRequest)
			return
		}

		// Create uploads directory
		uploadDir := "./uploads/articles"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			writeJSONError(w, "Gagal membuat direktori", http.StatusInternalServerError)
			return
		}

		// Generate filename
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		filePath := filepath.Join(uploadDir, filename)

		// Save file
		dst, err := os.Create(filePath)
		if err != nil {
			writeJSONError(w, "Gagal menyimpan file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			writeJSONError(w, "Gagal menyimpan file", http.StatusInternalServerError)
			return
		}

		// Return path
		relativePath := "/uploads/articles/" + filename
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"path":    relativePath,
		})
	}
}
