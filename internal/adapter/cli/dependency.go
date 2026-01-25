//go:generate mockgen -source=dependency.go -destination=dependency_mock.go -package=cli
package cli

import (
	"context"

	"github.com/jian-hua-he/reference-app/internal/entity"
)

type NoteApp interface {
	List(ctx context.Context) ([]entity.Note, error)
	Create(ctx context.Context, text string) (*entity.Note, error)
	Delete(ctx context.Context, id string) error
}
