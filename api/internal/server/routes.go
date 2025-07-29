package server

import (
	"net/http"

	"news-portal/api/internal/auth"

	"github.com/gorilla/mux"
)

// setupRoutes configures all routes on the provided router
func (s *Server) setupRoutes(router *mux.Router) {
	// Basic health and utility endpoints
	s.registerBasicRoutes(router)

	// Serve static files
	s.registerStaticRoutes(router)

	// Authentication routes (public)
	s.registerAuthRoutes(router)

	// Public API routes (no authentication required)
	s.registerPublicAPIRoutes(router)

	// Protected API routes (authentication required)
	s.registerProtectedAPIRoutes(router)
}

// registerBasicRoutes registers basic health and utility endpoints
func (s *Server) registerBasicRoutes(router *mux.Router) {
	router.HandleFunc("/", s.welcome).Methods("GET")
	router.HandleFunc("/ping", s.ping).Methods("GET")
	router.HandleFunc("/health", s.healthCheck).Methods("GET")
	router.HandleFunc("/db-test", s.dbTest).Methods("GET")
}

// registerStaticRoutes registers static file serving
func (s *Server) registerStaticRoutes(router *mux.Router) {
	// Serve static files (uploads, images, etc.)
	fs := http.FileServer(http.Dir("./uploads/"))
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", fs))
}

// registerAuthRoutes registers authentication routes
func (s *Server) registerAuthRoutes(router *mux.Router) {
	// Auth routes will be implemented in login.go
	s.RegisterAuthRoutes(router)
}

// registerPublicAPIRoutes registers public API routes (no auth required)
func (s *Server) registerPublicAPIRoutes(router *mux.Router) {
	publicAPI := router.PathPrefix("/api/v1").Subrouter()

	// Optional auth middleware for better user experience
	publicAPI.Use(s.jwtManager.OptionalAuthMiddleware)

	// Register public routes from other files
	s.RegisterPublicArticleRoutes(publicAPI)
	s.RegisterPublicTagRoutes(publicAPI)
	s.RegisterPublicCategoryRoutes(publicAPI)
	s.RegisterPublicCommentRoutes(publicAPI)
}

// registerProtectedAPIRoutes registers protected API routes (auth required)
func (s *Server) registerProtectedAPIRoutes(router *mux.Router) {
	protectedAPI := router.PathPrefix("/api/v1").Subrouter()
	protectedAPI.Use(s.jwtManager.AuthMiddleware)

	// User-level protected routes
	s.registerUserRoutes(protectedAPI)

	// Editor-level routes (editor + admin)
	s.registerEditorRoutes(protectedAPI)

	// Admin-only routes
	s.registerAdminRoutes(protectedAPI)
}

// registerUserRoutes registers authenticated user routes
func (s *Server) registerUserRoutes(router *mux.Router) {
	// User profile and personal data
	s.RegisterUserProtectedRoutes(router)

	// User interactions with content
	s.RegisterProtectedArticleRoutes(router)
	s.RegisterProtectedCommentRoutes(router)
}

// registerEditorRoutes registers editor and admin content management routes
func (s *Server) registerEditorRoutes(router *mux.Router) {
	editorAPI := router.PathPrefix("").Subrouter()
	editorAPI.Use(auth.EditorOrAdmin)

	// Content management routes for editors
	s.RegisterEditorRoutes(editorAPI)
}

// registerAdminRoutes registers admin-only routes
func (s *Server) registerAdminRoutes(router *mux.Router) {
	adminAPI := router.PathPrefix("/admin").Subrouter()
	adminAPI.Use(auth.AdminOnly)

	// Admin management routes
	s.RegisterAdminRoutes(adminAPI)
}

// Placeholder route registration functions
// These will be implemented in their respective files


func (s *Server) RegisterPublicArticleRoutes(router *mux.Router) {
	// Will be implemented in articles.go
}

func (s *Server) RegisterPublicTagRoutes(router *mux.Router) {
	// Will be implemented in tags.go
}

func (s *Server) RegisterPublicCategoryRoutes(router *mux.Router) {
	// Will be implemented in categories.go
}

func (s *Server) RegisterPublicCommentRoutes(router *mux.Router) {
	// Will be implemented in comments.go
}

func (s *Server) RegisterUserProtectedRoutes(router *mux.Router) {
	// Will be implemented in users.go or profile.go
}

func (s *Server) RegisterProtectedArticleRoutes(router *mux.Router) {
	// Will be implemented in articles.go
}

func (s *Server) RegisterProtectedCommentRoutes(router *mux.Router) {
	// Will be implemented in comments.go
}

func (s *Server) RegisterEditorRoutes(router *mux.Router) {
	// Will be implemented across multiple files (articles.go, tags.go, etc.)
}

func (s *Server) RegisterAdminRoutes(router *mux.Router) {
	// Will be implemented in admin.go or across multiple files
}
