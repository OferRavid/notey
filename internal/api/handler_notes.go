package api

import (
	"encoding/json"
	"net/http"

	"github.com/OferRavid/notey/internal/auth"
	"github.com/OferRavid/notey/internal/database"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (cfg *ApiConfig) handlerCreateNote(c echo.Context) error {
	bearerToken, err := auth.GetBearerToken(c.Request().Header)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Missing token in Authorization header")

	}

	type parameters struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	decoder := json.NewDecoder(c.Request().Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Couldn't decode parameters")

	}

	user_id, err := auth.ValidateJWT(bearerToken, cfg.Secret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Couldn't validate token")

	}

	note, err := cfg.DbQueries.CreateNote(c.Request().Context(), database.CreateNoteParams{
		Title:   params.Title,
		Content: params.Content,
		UserID:  user_id,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Couldn't create note")

	}

	return c.JSON(http.StatusCreated, Note{
		ID:        note.ID,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
		Title:     note.Title,
		Content:   note.Content,
		UserID:    note.UserID,
	})
}

func (cfg *ApiConfig) handlerRetrieveNotes(c echo.Context) error {
	author_id := c.Request().URL.Query().Get("author_id")
	// sortType := c.Request().URL.Query().Get("sort")
	dbNotes, err := cfg.DbQueries.GetNotes(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Couldn't retrieve notes")

	}

	notes := []Note{}
	for _, dbNote := range dbNotes {
		if author_id != "" && author_id != dbNote.UserID.String() {
			continue
		}
		notes = append(notes, Note{
			ID:        dbNote.ID,
			CreatedAt: dbNote.CreatedAt,
			UpdatedAt: dbNote.UpdatedAt,
			UserID:    dbNote.UserID,
			Title:     dbNote.Title,
			Content:   dbNote.Content,
		})
	}

	// if sortType == "desc" {
	// 	sort.Slice(notes, func(i, j int) bool { return notes[i].CreatedAt.After(notes[j].CreatedAt) })
	// }

	return c.JSON(http.StatusOK, notes)
}

func (cfg *ApiConfig) handlerGetNoteByID(c echo.Context) error {
	noteID, err := uuid.Parse(c.Request().PathValue("noteID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Failed to parse noteID")

	}
	note, err := cfg.DbQueries.GetNoteByID(c.Request().Context(), noteID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Couldn't retrieve note")

	}

	return c.JSON(http.StatusOK, Note{
		ID:        note.ID,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
		Title:     note.Title,
		Content:   note.Content,
		UserID:    note.UserID,
	})
}

// func (cfg *ApiConfig) handlerUpdateNote(c echo.Context) error

func (cfg *ApiConfig) handlerDeleteNote(c echo.Context) error {
	bearerToken, err := auth.GetBearerToken(c.Request().Header)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Missing token in Authorization header")

	}

	user_id, err := auth.ValidateJWT(bearerToken, cfg.Secret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Couldn't validate token")

	}

	noteID, err := uuid.Parse(c.Request().PathValue("noteID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Failed to parse noteID")

	}
	note, err := cfg.DbQueries.GetNoteByID(c.Request().Context(), noteID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Couldn't find note with the given ID")

	}

	if note.UserID != user_id {
		return c.JSON(http.StatusForbidden, "Unauthorized to delete note")

	}

	err = cfg.DbQueries.DeleteNote(c.Request().Context(), note.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to delete note")

	}

	return c.JSON(http.StatusNoContent, nil)
}
