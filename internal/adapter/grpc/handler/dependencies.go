//go:generate mockgen -source=dependencies.go -destination=dependencies_mock_test.go -package=handler_test
package handler

import (
	"context"

	"github.com/jian-hua-he/reference-app/internal/entity"
)

type NoteApp interface {
	List(ctx context.Context) ([]entity.Note, error)
	Create(ctx context.Context, text string) (*entity.Note, error)
	Delete(ctx context.Context, id string) error
}
