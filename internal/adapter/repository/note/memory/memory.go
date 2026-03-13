package memory

import (
	"context"
	"time"

	"github.com/jian-hua-he/reference-app/internal/entity"
	"github.com/jian-hua-he/reference-app/internal/adapter/repository"
)

type Repo struct {
	uuidFunc func() string
	nowFunc  func() time.Time
	notes    map[string]Note
}

func NewRepo(uuidFunc func() string, nowFunc func() time.Time) *Repo {
	return &Repo{
		uuidFunc: uuidFunc,
		nowFunc:  nowFunc,
		notes:    make(map[string]Note),
	}
}

func NewRepoWithNotes(
	uuidFunc func() string,
	nowFunc func() time.Time,
	notes map[string]Note,
) *Repo {
	if notes == nil {
		notes = make(map[string]Note)
	}

	return &Repo{
		uuidFunc: uuidFunc,
		nowFunc:  nowFunc,
		notes:    notes,
	}
}

func (r *Repo) Create(ctx context.Context, text string) (*entity.Note, error) {
	id := r.uuidFunc()
	createdAt := r.nowFunc()

	note := Note{
		Text:      text,
		CreatedAt: createdAt,
	}

	r.notes[id] = note

	return &entity.Note{
		ID:        id,
		Text:      text,
		CreatedAt: createdAt,
	}, nil
}

func (r *Repo) List(ctx context.Context) ([]entity.Note, error) {
	var result []entity.Note

	if len(r.notes) == 0 {
		return []entity.Note{}, nil
	}

	for id, note := range r.notes {
		result = append(result, entity.Note{
			ID:        id,
			Text:      note.Text,
			CreatedAt: note.CreatedAt,
		})
	}

	return result, nil
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	if _, exists := r.notes[id]; !exists {
		return repository.ErrNotFound
	}

	delete(r.notes, id)

	return nil
}
