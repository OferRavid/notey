package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"

	"github.com/OferRavid/notey/internal/database"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type NoteParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Handles creation of a new note for user.
func (cfg *ApiConfig) handlerCreateNote(c echo.Context) error {

	decoder := json.NewDecoder(c.Request().Body)
	params := NoteParams{}
	err := decoder.Decode(&params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Couldn't decode parameters"})

	}

	user_id := c.Get("user_id").(uuid.UUID)

	note, err := cfg.DbQueries.CreateNote(c.Request().Context(), database.CreateNoteParams{
		Title:   params.Title,
		Content: params.Content,
		UserID:  user_id,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Couldn't create note"})

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

// Handles getting all the user's notes.
func (cfg *ApiConfig) handlerRetrieveNotes(c echo.Context) error {
	notes := []Note{}
	user_id := c.Get("user_id").(uuid.UUID)
	dbNotes, err := cfg.DbQueries.GetNotesByUserID(c.Request().Context(), user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNoContent, notes)
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Couldn't retrieve notes"})
	}

	for _, dbNote := range dbNotes {
		notes = append(notes, Note{
			ID:        dbNote.ID,
			CreatedAt: dbNote.CreatedAt,
			UpdatedAt: dbNote.UpdatedAt,
			UserID:    dbNote.UserID,
			Title:     dbNote.Title,
			Content:   dbNote.Content,
		})
	}

	sort.Slice(notes, func(i, j int) bool { return notes[i].CreatedAt.Before(notes[j].CreatedAt) })

	return c.JSON(http.StatusOK, notes)
}

// Handles getting one note using it's ID.
func (cfg *ApiConfig) handlerGetNoteByID(c echo.Context) error {
	noteID, err := uuid.Parse(c.Param("noteID"))
	user_id := c.Get("user_id").(uuid.UUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"Error": "Failed to parse noteID"})

	}
	note, err := cfg.DbQueries.GetNoteByID(c.Request().Context(), noteID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"Error": "Couldn't retrieve note"})

	}

	if user_id != note.UserID {
		return c.JSON(http.StatusForbidden, echo.Map{"Error": "Unauthorized to view note"})
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

// Handles updating content and/or title of a note.
func (cfg *ApiConfig) handlerUpdateNote(c echo.Context) error {

	user_id := c.Get("user_id").(uuid.UUID)

	noteID, err := uuid.Parse(c.Param("noteID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"Error": "Failed to parse noteID"})

	}
	note, err := cfg.DbQueries.GetNoteByID(c.Request().Context(), noteID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"Error": "Couldn't find note with the given ID"})

	}

	if note.UserID != user_id {
		return c.JSON(http.StatusForbidden, echo.Map{"Error": "Unauthorized to edit note"})
	}

	noteParams := NoteParams{}
	decoder := json.NewDecoder(c.Request().Body)
	err = decoder.Decode(&noteParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "couldn't decode parameters"})
	}

	updatedNote, err := cfg.DbQueries.UpdateNote(
		c.Request().Context(),
		database.UpdateNoteParams{
			Title:   noteParams.Title,
			Content: noteParams.Content,
			ID:      noteID,
		},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Couldn't update user"})
	}

	return c.JSON(http.StatusOK,
		Note{
			ID:        updatedNote.ID,
			CreatedAt: updatedNote.CreatedAt,
			UpdatedAt: updatedNote.UpdatedAt,
			Title:     updatedNote.Title,
			Content:   updatedNote.Content,
		},
	)
}

// Handles deletion of a note.
func (cfg *ApiConfig) handlerDeleteNote(c echo.Context) error {
	user_id := c.Get("user_id").(uuid.UUID)

	noteIDStr := c.Param("noteID")

	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"Error": "Failed to parse noteID"})
	}
	note, err := cfg.DbQueries.GetNoteByID(c.Request().Context(), noteID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"Error": "Couldn't find note with the given ID"})
	}

	if note.UserID != user_id {
		return c.JSON(http.StatusForbidden, echo.Map{"Error": "Unauthorized to delete note"})
	}

	err = cfg.DbQueries.DeleteNote(c.Request().Context(), note.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Failed to delete note"})
	}

	return c.JSON(http.StatusNoContent, nil)
}
