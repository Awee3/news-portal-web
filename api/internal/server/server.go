package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"news-portal-web/api/internal/auth"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Server holds dependencies for HTTP handlers
type Server struct {
	db         *sql.DB
	jwtManager *auth.JWTManager
}

// NewServer creates a new server instance
// Ganti fungsi NewServer menjadi:
func NewServer(db *sql.DB, secretKey string) *Server {
	jwtManager := auth.NewJWTManager(secretKey)

	return &Server{
		db:         db,
		jwtManager: jwtManager,
	}
}

// GetDB returns the database connection
func (s *Server) GetDB() *sql.DB {
	return s.db
}

// GetJWTManager returns the JWT manager
func (s *Server) GetJWTManager() *auth.JWTManager {
	return s.jwtManager
}

// Start starts the server with CORS enabled
// Ganti fungsi Start menjadi:
func (s *Server) Start(addr string) error {
	router := s.SetupRoutes()

	// Debug: log semua routes yang terdaftar
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		if path != "" {
			log.Printf("📍 Route: %v %s", methods, path)
		}
		return nil
	})

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Printf("🚀 Server listening on %s", addr)
	return http.ListenAndServe(addr, handler)
}

// ========================================
// BASIC HANDLERS
// ========================================

func (s *Server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Check database connection
		if err := s.db.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"status":   "unhealthy",
				"service":  "news-portal-api",
				"database": "disconnected",
				"error":    err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":   "healthy",
			"service":  "news-portal-api",
			"database": "connected",
		})
	}
}

func (s *Server) handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "pong",
		})
	}
}

func (s *Server) handleWelcome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Welcome to News Portal API",
			"version": "1.0.0",
		})
	}
}

func (s *Server) handleDBTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Test simple query
		var result int
		err := s.db.QueryRow("SELECT 1").Scan(&result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "ok",
			"database":   "connected",
			"test_query": result,
		})
	}
}
