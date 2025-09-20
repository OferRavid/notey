package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/OferRavid/notey/internal/api"
	"github.com/OferRavid/notey/internal/database"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

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
	dbQueries := database.New(db)
	apiCfg := &api.ApiConfig{
		DbQueries: dbQueries,
		Platform:  platform,
		Secret:    jwtSecret,
	}

	e := echo.New()
	apiCfg.RegisterRoutes(e)
	e.Static("/", filepathRoot)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	e.Logger.Fatal(e.Start(":" + port))
}
