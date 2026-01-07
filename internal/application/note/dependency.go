//go:generate mockgen -source=dependency.go -destination=dependency_mock.go -package=note
package note

import (
	"context"

	"github.com/jian-hua-he/ddd_notes/internal/entity"
)

type NoteRepository interface {
	List(ctx context.Context) ([]entity.Note, error)
	Create(ctx context.Context, text string) (*entity.Note, error)
	Delete(ctx context.Context, id string) error
}
