package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	noteApp NoteApp
}

func NewHandler(noteApp NoteApp) *Handler {
	return &Handler{
		noteApp: noteApp,
	}
}

// GetNotes
//
// @Summary List Notes
// @Description Get a list of all notes
// @Tags v1
// @Accept json
// @Produce json
// @Success 200 {object} GetNotesResponse
// @Failure 500 {object} GetNotesResponse
// @Router /v1/notes [GET]
func (h *Handler) GetNotes(c echo.Context) error {
	notes, err := h.noteApp.List(c.Request().Context())
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			GetNotesResponse{
				Message: "failed to list notes",
			},
		)
	}

	payload := make([]Note, len(notes))
	for i, note := range notes {
		payload[i] = Note{
			ID:        note.ID,
			Text:      note.Text,
			CreatedAt: note.CreatedAt,
		}
	}

	return c.JSON(
		http.StatusOK,
		GetNotesResponse{
			Message: "successful",
			Payload: payload,
		},
	)
}

// PostNote
//
// @Summary Create Note
// @Description Create a new note
// @Tags v1
// @Accept json
// @Produce json
// @Param request body PostNoteRequest true "Note text"
// @Success 200 {object} PostNoteResponse
// @Failure 400 {object} PostNoteResponse
// @Failure 500 {object} PostNoteResponse
// @Router /v1/notes [POST]
func (h *Handler) PostNote(c echo.Context) error {
	var req PostNoteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			PostNoteResponse{
				Message: "invalid request body",
			},
		)
	}

	note, err := h.noteApp.Create(c.Request().Context(), req.Text)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			PostNoteResponse{
				Message: "failed to create note",
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		PostNoteResponse{
			Message: "successful",
			Payload: &Note{
				ID:        note.ID,
				Text:      note.Text,
				CreatedAt: note.CreatedAt,
			},
		},
	)
}

// DeleteNote
//
// @Summary Delete Note
// @Description Delete a note by ID
// @Tags v1
// @Accept json
// @Produce json
// @Param note_id path string true "Note ID"
// @Success 200 {object} DeleteNoteResponse
// @Failure 500 {object} DeleteNoteResponse
// @Router /v1/notes/{note_id} [DELETE]
func (h *Handler) DeleteNote(c echo.Context) error {
	noteID := c.Param("note_id")

	err := h.noteApp.Delete(c.Request().Context(), noteID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			DeleteNoteResponse{
				Message: "failed to delete note",
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		DeleteNoteResponse{
			Message: "successful",
		},
	)
}
