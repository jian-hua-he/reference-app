package web

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

func (h *Handler) ListNotes(c echo.Context) error {
	notes, err := h.noteApp.List(c.Request().Context())
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			ErrResponse{Message: "failed to list notes"},
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
		ListNoteResponse{
			Message: "successful",
			Payload: payload,
		},
	)
}
