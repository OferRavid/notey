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

	"github.com/OferRavid/notey/internal/auth"
	"github.com/OferRavid/notey/internal/database"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// Test the handlerLogout function
func TestHandlerLogout(t *testing.T) {
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

	if err := cfg.DbQueries.DeleteUsers(context.Background()); err != nil {
		t.Fatalf("failed to delete users: %s", err)
	}

	_, err = cfg.DbQueries.CreateUser(context.Background(), database.CreateUserParams{
		Username:       "testuser",
		Email:          "test@example.com",
		HashedPassword: hashedPassword,
	})
	if err != nil {
		t.Fatalf("failed to insert test user: %v", err)
	}

	// Test case: successful logout
	t.Run("successful logout", func(t *testing.T) {
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
		if err != nil {
			t.Fatalf("failed to login test user: %v", err)
		}

		var loginResponse struct {
			ID    uuid.UUID `json:"id"`
			Token string    `json:"token"`
		}
		json.Unmarshal(rec.Body.Bytes(), &loginResponse)
		jwtToken := loginResponse.Token
		userID := loginResponse.ID

		count, err := cfg.DbQueries.CheckRecordExists(c.Request().Context(), userID)
		assert.NoError(t, err)
		assert.Greater(t, count, int64(0))

		req = httptest.NewRequest(http.MethodDelete, "/api/logout", &bytes.Buffer{})
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.Set("user_id", userID)

		err = cfg.handlerLogout(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusNoContent {
			fmt.Fprintf(os.Stderr, "response returned error: %s", rec.Body.String())
		}

		assert.Equal(t, http.StatusNoContent, rec.Code)
		count, err = cfg.DbQueries.CheckRecordExists(c.Request().Context(), userID)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}
