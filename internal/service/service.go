package service

import (
	"context"
	"errors"

	"github.com/jian-hua-he/ddd_notes/internal/repository"
)

type NoteService struct {
	repo NoteRepository
}

func NewNoteService(repo NoteRepository) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) Create(ctx context.Context, text string) (*ServiceNote, error) {
	n, err := s.repo.Create(ctx, text)
	if err != nil {
		return nil, err
	}

	return &ServiceNote{
		ID:        n.ID,
		Text:      n.Text,
		CreatedAt: n.CreatedAt,
	}, nil
}

func (s *NoteService) List(ctx context.Context) ([]ServiceNote, error) {
	ns, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]ServiceNote, 0, len(ns))
	for _, n := range ns {
		out = append(out, ServiceNote{
			ID:        n.ID,
			Text:      n.Text,
			CreatedAt: n.CreatedAt,
		})
	}

	return out, nil
}

func (s *NoteService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err == nil {
		return nil
	}
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}

	return err
}
