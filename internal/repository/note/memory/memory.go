package memory

import (
	"context"
	"errors"

	"github.com/jian-hua-he/ddd_notes/internal/repository/note"
)

type Repo struct{}

func (r *Repo) Create(ctx context.Context, text string) (*note.RepoNote, error) {
	return nil, errors.New("not implemented")
}

func (r *Repo) List(ctx context.Context) ([]note.RepoNote, error) {
	return nil, errors.New("not implemented")
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}
