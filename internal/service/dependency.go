//go:generate mockgen -source=dependency.go -destination=dependency_mock.go -package=service
package service

import (
	"context"

	"github.com/jian-hua-he/ddd_notes/internal/repository"
)

type NoteRepository interface {
	Create(ctx context.Context, text string) (*repository.RepoNote, error)
	List(ctx context.Context) ([]repository.RepoNote, error)
	Delete(ctx context.Context, id string) error
}
