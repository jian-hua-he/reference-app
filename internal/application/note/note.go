package note

import (
	"context"
	"errors"
	"fmt"

	"github.com/jian-hua-he/ddd_notes/internal/application"
	"github.com/jian-hua-he/ddd_notes/internal/entity"
	"github.com/jian-hua-he/ddd_notes/internal/repository"
)

type NoteApp struct {
	repo NoteRepository
}

func NewNoteApp(repo NoteRepository) *NoteApp {
	return &NoteApp{repo: repo}
}

func (s *NoteApp) List(ctx context.Context) ([]entity.Note, error) {
	notes, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list all notes: %w", err)
	}

	result := make([]entity.Note, 0, len(notes))
	for _, n := range notes {
		result = append(result, entity.Note{
			ID:        n.ID,
			Text:      n.Text,
			CreatedAt: n.CreatedAt,
		})
	}

	return result, nil
}

func (s *NoteApp) Create(ctx context.Context, text string) (*entity.Note, error) {
	n, err := s.repo.Create(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("failed to create a note: %w", err)
	}

	return &entity.Note{
		ID:        n.ID,
		Text:      n.Text,
		CreatedAt: n.CreatedAt,
	}, nil
}

func (s *NoteApp) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return application.ErrNotFound
	}

	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	return nil
}
