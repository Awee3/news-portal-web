package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"news-portal-web/api/internal/server"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file
	_ = godotenv.Load()

	// Get environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Build from components
		dbHost := getEnvWithDefault("DB_HOST", "localhost")
		dbPort := getEnvWithDefault("DB_PORT", "5432")
		dbUser := getEnvWithDefault("DB_USER", "postgres")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := getEnvWithDefault("DB_NAME", "Winnicode")

		if dbPassword == "" {
			log.Fatal("❌ Database password required: set DB_PASSWORD or DATABASE_URL")
		}

		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			dbUser, dbPassword, dbHost, dbPort, dbName)
	}

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		log.Fatal("❌ JWT_SECRET_KEY environment variable is required")
	}

	port := getEnvWithDefault("PORT", "8080")
	env := getEnvWithDefault("ENV", "development")

	// Connect to database
	log.Printf("🔌 Connecting to database...")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("❌ Failed to ping database:", err)
	}
	log.Printf("✅ Database connected successfully")

	

	// Create server instance
	srv := server.NewServer(db, jwtSecret)

	// Start server
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("🌍 Environment: %s", env)
	log.Printf("📊 Health: http://localhost:%s/health", port)
	log.Printf("🏓 Ping: http://localhost:%s/ping", port)
	log.Printf("🗄️  DB Test: http://localhost:%s/db-test", port)

	if err := srv.Start(":" + port); err != nil {
		log.Fatal("❌ Server failed to start:", err)
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
