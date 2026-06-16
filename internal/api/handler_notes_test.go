package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
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

// Test the handlerCreateNote function
func TestHandlerCreateNote(t *testing.T) {
	e := echo.New()

	// Setup test database
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("failed to setup test DB: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Create Queries instance
	queries := database.New(db)
	cfg := &ApiConfig{
		DbQueries: queries,
		Secret:    "secret",
	}

	// Define cleanup function
	t.Cleanup(func() {
		cfg.DbQueries.DeleteUsers(context.Background())
		db.Close()
	})

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

	body := map[string]string{
		"username": "testuser",
		"password": "password",
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = cfg.handlerLogin(c)
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

	tests := []struct {
		name           string
		title          string
		content        string
		expectedStatus int
	}{
		{
			name:           "Note created successfully",
			title:          "Test note creation title",
			content:        "Test note creation content",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Empty title and content",
			title:          "",
			content:        "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body = map[string]string{
				"title":   tt.title,
				"content": tt.content,
			}
			bodyBytes, _ := json.Marshal(body)
			req := httptest.NewRequest(http.MethodPost, "/api/notes", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Authorization", "Bearer"+jwtToken)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user_id", userID)

			err = cfg.handlerCreateNote(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusCreated {
				var response Note
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Equal(t, tt.title, response.Title)
				assert.Equal(t, tt.content, response.Content)
				assert.Equal(t, userID, response.UserID)
			}

		})
	}

}
