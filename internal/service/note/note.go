package note

import (
	"context"
	"errors"
	"fmt"

	"github.com/jian-hua-he/ddd_notes/internal/domain"
	"github.com/jian-hua-he/ddd_notes/internal/repository"
	"github.com/jian-hua-he/ddd_notes/internal/service"
)

type NoteService struct {
	repo NoteRepository
}

func NewNoteService(repo NoteRepository) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) Create(ctx context.Context, text string) (*domain.Note, error) {
	n, err := s.repo.Create(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("failed to create a note: %w", err)
	}

	return &domain.Note{
		ID:        n.ID,
		Text:      n.Text,
		CreatedAt: n.CreatedAt,
	}, nil
}

func (s *NoteService) List(ctx context.Context) ([]domain.Note, error) {
	notes, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list all notes: %w", err)
	}

	result := make([]domain.Note, 0, len(notes))
	for _, n := range notes {
		result = append(result, domain.Note{
			ID:        n.ID,
			Text:      n.Text,
			CreatedAt: n.CreatedAt,
		})
	}

	return result, nil
}

func (s *NoteService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return service.ErrNotFound
	}

	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	return nil
}
