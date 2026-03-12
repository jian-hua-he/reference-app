package postgres

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/google/uuid"

	"github.com/jian-hua-he/reference-app/internal/entity"
	"github.com/jian-hua-he/reference-app/internal/repository"
)

var dialect = goqu.Dialect("postgres")

const tableName = "notes"

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, text string) (*entity.Note, error) {
	id := uuid.New().String()

	query, args, err := dialect.
		Insert(tableName).
		Cols("id", "text").
		Vals(goqu.Vals{id, text}).
		Returning("id", "text", "created_at").
		ToSQL()
	if err != nil {
		return nil, err
	}

	var note entity.Note
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&note.ID, &note.Text, &note.CreatedAt); err != nil {
		return nil, err
	}

	return &note, nil
}

func (r *Repo) List(ctx context.Context) ([]entity.Note, error) {
	query, args, err := dialect.
		From(tableName).
		Select("id", "text", "created_at").
		Order(goqu.C("created_at").Asc()).
		ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
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
	query, args, err := dialect.
		Delete(tableName).
		Where(goqu.C("id").Eq(id)).
		ToSQL()
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, query, args...)
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
