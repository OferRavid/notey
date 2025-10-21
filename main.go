package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/OferRavid/notey/internal/api"
	"github.com/OferRavid/notey/internal/database"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	const filepathRoot = "."
	const staticDir = "static"
	const port = "8080"

	// Use environment variables to set up server's environment and saving them in a struct
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open db: %s\n", err)
	}
	defer db.Close()

	// Use db.Ping() to confirm the connection is valid
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	dbQueries := database.New(db)

	cfg := &api.ApiConfig{
		FileserverHits: atomic.Int32{},
		DbQueries:      dbQueries,
		FilepathRoot:   filepathRoot,
		StaticDir:      staticDir,
		Platform:       platform,
		Secret:         jwtSecret,
	}

	// Set up the server
	e := echo.New()

	go func() {
		// Create a ticker that ticks every hour.
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			cfg.DbQueries.ClearRevokedTokens(context.Background())
		}
	}()

	e.Static("/static", "static")

	cfg.RegisterRoutes(e)

	// Create a group for routes starting with /app
	appGroup := e.Group("/app")

	// Serve static files and strip the /app prefix
	// Use the custom handler to serve static files
	appGroup.GET("/*", cfg.ServeStaticFiles)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	e.Logger.Fatal(e.Start(":" + port))
}
