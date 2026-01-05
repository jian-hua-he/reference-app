//go:generate mockgen -source=dependency.go -destination=dependency_mock.go -package=service
package service

import (
	"context"

	"github.com/jian-hua-he/ddd_notes/internal/repository/note"
)

type NoteRepository interface {
	List(ctx context.Context) ([]note.RepoNote, error)
	Create(ctx context.Context, text string) (*note.RepoNote, error)
	Delete(ctx context.Context, id string) error
}
