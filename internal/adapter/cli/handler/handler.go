package handler

import (
	"context"
	"fmt"
	"io"
)

type Handler struct {
	noteApp NoteApp
	out     io.Writer
}

func NewHandler(noteApp NoteApp, out io.Writer) *Handler {
	return &Handler{
		noteApp: noteApp,
		out:     out,
	}
}

func (h *Handler) ListNotes(ctx context.Context) error {
	notes, err := h.noteApp.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	if len(notes) == 0 {
		fmt.Fprintln(h.out, "No notes found.")
		return nil
	}

	for _, n := range notes {
		fmt.Fprintf(h.out, "ID: %s\nText: %s\nCreated At: %s\n\n", n.ID, n.Text, n.CreatedAt)
	}

	return nil
}

func (h *Handler) CreateNote(ctx context.Context, text string) error {
	note, err := h.noteApp.Create(ctx, text)
	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	fmt.Fprintln(h.out, "Note created successfully.")
	fmt.Fprintf(h.out, "ID: %s\nText: %s\nCreated At: %s\n", note.ID, note.Text, note.CreatedAt)

	return nil
}

func (h *Handler) DeleteNote(ctx context.Context, id string) error {
	err := h.noteApp.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	fmt.Fprintln(h.out, "Note deleted successfully.")

	return nil
}
