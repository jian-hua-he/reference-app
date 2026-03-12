package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/jian-hua-he/reference-app/internal/entity"
	"github.com/jian-hua-he/reference-app/internal/repository"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, text string) (*entity.Note, error) {
	id := uuid.New().String()

	var note entity.Note
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO notes (id, text) VALUES ($1, $2) RETURNING id, text, created_at",
		id, text,
	).Scan(&note.ID, &note.Text, &note.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (r *Repo) List(ctx context.Context) ([]entity.Note, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, text, created_at FROM notes ORDER BY created_at")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []entity.Note
	for rows.Next() {
		var n entity.Note
		if err := rows.Scan(&n.ID, &n.Text, &n.CreatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if notes == nil {
		return []entity.Note{}, nil
	}

	return notes, nil
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM notes WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}
