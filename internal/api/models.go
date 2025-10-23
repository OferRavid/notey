package api

import (
	"sync/atomic"
	"time"

	"github.com/OferRavid/notey/internal/database"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
)

type ApiConfig struct {
	FileserverHits  atomic.Int32
	PageVisitsGauge prometheus.Gauge
	DbQueries       *database.Queries
	StaticDir       string
	Platform        string
	Secret          string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Note struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    uuid.UUID `json:"user_id"`
}
