package server

import (
	"net/http"

	"news-portal-web/api/internal/auth"

	"github.com/gorilla/mux"
)

// SetupRoutes configures and returns the router with all API routes
func (s *Server) SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// ========================================
	// API v1 ROUTER
	// ========================================
	api := r.PathPrefix("/api/v1").Subrouter()

	// ========================================
    // BASIC ROUTES (tanpa prefix)
    // ========================================
    r.HandleFunc("/", s.handleWelcome()).Methods("GET")
    r.HandleFunc("/health", s.handleHealth()).Methods("GET")
    r.HandleFunc("/ping", s.handlePing()).Methods("GET")
    r.HandleFunc("/db-test", s.handleDBTest()).Methods("GET")
	
	// ========================================
	// PUBLIC ROUTES (no auth required)
	// ========================================
	public := api.NewRoute().Subrouter()

	// Articles - public endpoints
	s.RegisterPublicArticleRoutes(public)

	// Categories - public endpoints
	s.RegisterPublicCategoryRoutes(public)

	// Tags - public endpoints
	s.RegisterPublicTagRoutes(public)

	// Comments - public endpoints (get & create)
	s.RegisterPublicCommentRoutes(public)

	// Auth routes (login, register, refresh)
	public.HandleFunc("/auth/login", s.handleLogin()).Methods("POST")
	public.HandleFunc("/auth/register", s.handleRegister()).Methods("POST")
	public.HandleFunc("/auth/refresh", s.handleRefreshToken()).Methods("POST")

	// ========================================
	// AUTHENTICATED USER ROUTES
	// ========================================
	authenticated := api.NewRoute().Subrouter()
	authenticated.Use(auth.AuthMiddleware(s.GetJWTManager()))

	// Auth - logout (requires token)
	authenticated.HandleFunc("/auth/logout", s.handleLogout()).Methods("POST")

	// User profile routes
	s.RegisterUserRoutes(authenticated)

	// User comment routes
	s.RegisterUserCommentRoutes(authenticated)

	// ========================================
	// EDITOR ROUTES (editor + admin only)
	// ========================================
	editor := api.PathPrefix("/editor").Subrouter()
	editor.Use(auth.AuthMiddleware(s.GetJWTManager()))
	editor.Use(auth.EditorOrAdmin)

	// Editor article management
	s.RegisterEditorArticleRoutes(editor)

	// ========================================
	// ADMIN ROUTES (admin only)
	// ========================================
	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(auth.AuthMiddleware(s.GetJWTManager()))
	admin.Use(auth.AdminOnly)

	// User management
	s.RegisterAdminUserRoutes(admin)

	// Category management
	s.RegisterAdminCategoryRoutes(admin)

	// Tag management
	s.RegisterAdminTagRoutes(admin)

	// Comment moderation
	s.RegisterAdminCommentRoutes(admin)

	// ========================================
	// STATIC FILES
	// ========================================
	r.PathPrefix("/uploads/").Handler(
		http.StripPrefix("/uploads/",
			http.FileServer(http.Dir("./uploads"))))

	return r
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

