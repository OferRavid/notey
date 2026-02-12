package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/OferRavid/notey/internal/auth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// Test the middleware function
func TestMiddleware(t *testing.T) {
	e := echo.New()
	cfg := &ApiConfig{
		Secret: "test_secret",
	}
	e.Use(cfg.Middleware())
	userID := uuid.New()
	validToken, _ := auth.MakeJWT(userID, cfg.Secret, time.Hour)
	invalidJWTToken, _ := auth.MakeJWT(userID, "invalid_test_secret", time.Hour)
	tests := []struct {
		name             string
		header           string
		expectedStatus   int
		expectedResponse map[string]interface{}
	}{

		{
			name:             "No Authorization Header",
			header:           "",
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: map[string]interface{}{"error": "invalid Authorization header format"},
		},

		{
			name:             "Invalid Token Format",
			header:           "Bearer ",
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: map[string]interface{}{"error": "invalid Authorization header format"},
		},

		{
			name:             "Valid Token but Invalid JWT",
			header:           "Bearer " + invalidJWTToken,
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: map[string]interface{}{"error": "invalid token"},
		},

		{
			name:             "Valid Token",
			header:           "Bearer " + validToken,
			expectedStatus:   http.StatusOK,
			expectedResponse: map[string]interface{}{},
		},
	}
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Call the middleware and then the handler
			if err := cfg.Middleware()(handler)(c); err != nil {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Equal(t, tt.expectedResponse, rec.Body)

			}
			if tt.name == "Valid Token" {
				fmt.Println("Check middleware worked as expected...")
				assert.Equal(t, userID, c.Get("user_id").(uuid.UUID))
			}

		})
	}
}
