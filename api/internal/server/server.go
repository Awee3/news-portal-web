package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"news-portal/api/internal/auth"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	db         *sql.DB
	jwtManager *auth.JWTManager
}

// NewServer creates a new server instance
func NewServer(db *sql.DB, secretKey string) *Server {
	// JWT token durations
	accessTokenTTL := 15 * time.Minute    // 15 minutes
	refreshTokenTTL := 7 * 24 * time.Hour // 7 days

	jwtManager := auth.NewJWTManager(secretKey, accessTokenTTL, refreshTokenTTL)

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

// RegisterRoutes returns configured router with all routes
func (s *Server) RegisterRoutes() http.Handler {
	router := mux.NewRouter()

	// Register all routes (implemented in routes.go)
	s.setupRoutes(router)

	return router
}

// Start starts the server with CORS enabled
func (s *Server) Start(addr string) error {
	router := s.RegisterRoutes()

	// Setup CORS for frontend
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001"}, // Next.js dev server
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	return http.ListenAndServe(addr, handler)
}

// Basic handlers
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "pong",
	})
}

func (s *Server) welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Welcome to News Portal API",
		"version": "1.0.0",
		"status":  "running",
	})
}

// Test koneksi database
func (s *Server) dbTest(w http.ResponseWriter, r *http.Request) {
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
