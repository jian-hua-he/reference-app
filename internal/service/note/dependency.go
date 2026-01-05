//go:generate mockgen -source=dependency.go -destination=dependency_mock.go -package=note
package note

import (
	"context"

	"github.com/jian-hua-he/ddd_notes/internal/domain"
)

type NoteRepository interface {
	List(ctx context.Context) ([]domain.Note, error)
	Create(ctx context.Context, text string) (*domain.Note, error)
	Delete(ctx context.Context, id string) error
}
