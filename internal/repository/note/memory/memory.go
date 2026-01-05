package memory

import (
	"context"
	"errors"

	"github.com/jian-hua-he/ddd_notes/internal/domain"
)

type Repo struct{}

func (r *Repo) Create(ctx context.Context, text string) (*domain.Note, error) {
	return nil, errors.New("not implemented")
}

func (r *Repo) List(ctx context.Context) ([]domain.Note, error) {
	return nil, errors.New("not implemented")
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}
