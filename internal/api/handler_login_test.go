package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/OferRavid/notey/internal/auth"
	"github.com/OferRavid/notey/internal/database"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// Test the handlerLogin function
func TestHandlerLogin(t *testing.T) {
	e := echo.New()

	// Setup test database
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("failed to setup test DB: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Create Queries instance
	queries := database.New(db)
	cfg := &ApiConfig{
		DbQueries: queries,
		Secret:    "secret",
	}

	// Insert test user into the database
	hashedPassword, err := auth.HashPassword("password")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	_, err = cfg.DbQueries.CreateUser(context.Background(), database.CreateUserParams{
		Username:       "testuser",
		Email:          "test@example.com",
		HashedPassword: hashedPassword,
	})
	if err != nil {
		t.Fatalf("failed to insert test user: %v", err)
	}

	// Test case: successful login
	t.Run("successful login", func(t *testing.T) {
		body := map[string]string{
			"username": "testuser",
			"password": "password", // plain text password
		}
		bodyBytes, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := cfg.handlerLogin(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			fmt.Fprintf(os.Stderr, "response returned error: %s", rec.Body.String())
		}

		assert.Equal(t, http.StatusOK, rec.Code)

		var response struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Username  string    `json:"username"`
			Email     string    `json:"email"`
			Token     string    `json:"token"`
		}
		fmt.Println("Body contains: ")
		fmt.Println(rec.Body.String())
		json.Unmarshal(rec.Body.Bytes(), &response)
		fmt.Println("response returned username: ")
		fmt.Println(response.Username)
		assert.Equal(t, "testuser", response.Username)
		assert.NotEmpty(t, response.Token) // Check that a token was returned
	})

	// Test case: user not found
	t.Run("user not found", func(t *testing.T) {
		body := map[string]string{
			"username": "wronguser",
			"password": "password",
		}
		bodyBytes, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := cfg.handlerLogin(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	// Test case: incorrect password
	t.Run("incorrect password", func(t *testing.T) {
		body := map[string]string{
			"username": "testuser",
			"password": "wrongpassword", // incorrect password
		}
		bodyBytes, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := cfg.handlerLogin(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	// Additional test cases can be added for other scenarios (e.g., JWT creation failure, etc.)
}
