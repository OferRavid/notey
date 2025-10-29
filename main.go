package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/OferRavid/notey/internal/api"
	"github.com/OferRavid/notey/internal/database"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"

	_ "github.com/lib/pq"
)

// The path where Docker secrets are mounted
const secretsDir = "/run/secrets"

func main() {
	const staticDir = "static"
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080" // Provide a default if not set
	}

	var pageVisitsGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "page_visits_gauge",
			Help: "Current number of visits to the app page.",
		},
	)

	databaseUser, err := readSecretFile("db_user")
	if err != nil {
		log.Fatalf("Failed to read database user secret: %v", err)
	}

	databasePassword, err := readSecretFile("db_password")
	if err != nil {
		log.Fatalf("Failed to read database password secret: %v", err)
	}

	jwtSecret, err := readSecretFile("jwt_key")
	if err != nil {
		log.Fatalf("Failed to read JWT secret: %v", err)
	}

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@db:5432/app_db?sslmode=disable",
		databaseUser,
		databasePassword,
	)

	platform := os.Getenv("PLATFORM")
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
		FileserverHits:  atomic.Int32{},
		PageVisitsGauge: pageVisitsGauge,
		DbQueries:       dbQueries,
		StaticDir:       staticDir,
		Platform:        platform,
		Secret:          jwtSecret,
	}

	// Set up the server
	e := echo.New()

	// Register your custom metric before using it
	prometheus.MustRegister(cfg.PageVisitsGauge)

	// Add middleware for collecting standard metrics
	e.Use(echoprometheus.NewMiddleware("myapp"))

	// Route to expose the combined metrics
	e.GET("/metrics", echoprometheus.NewHandler())

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

	log.Printf("Serving files on port: %s\n", appPort)
	if cfg.Platform == "dev" {
		log.Println("starting non-secure http page")
		e.Logger.Fatal(e.Start(":" + appPort))
	} else {
		e.Logger.Fatal(e.StartAutoTLS(":" + appPort))
	}

}

// readSecretFile reads a secret from its mounted file path.
func readSecretFile(secretName string) (string, error) {
	filePath := fmt.Sprintf("%s/%s", secretsDir, secretName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	// Trim any trailing newline characters
	return string(content), nil
}
