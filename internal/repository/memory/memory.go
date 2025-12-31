package memory

import (
	"context"
	"errors"

	"github.com/jian-hua-he/ddd_notes/internal/repository"
)

type Repo struct{}

func (r *Repo) Create(ctx context.Context, text string) (*repository.RepoNote, error) {
	return nil, errors.New("not implemented")
}

func (r *Repo) List(ctx context.Context) ([]repository.RepoNote, error) {
	return nil, errors.New("not implemented")
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}
