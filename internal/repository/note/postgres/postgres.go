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

const tableName = "notes"

type Repo struct {
	db *goqu.Database
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: goqu.Dialect("postgres").DB(db)}
}

func (r *Repo) Create(ctx context.Context, text string) (*entity.Note, error) {
	id := uuid.New().String()

	var note entity.Note
	found, err := r.db.
		Insert(tableName).
		Cols("id", "text").
		Vals(goqu.Vals{id, text}).
		Returning("id", "text", "created_at").
		Executor().
		ScanStructContext(ctx, &note)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, repository.ErrNotFound
	}

	return &note, nil
}

func (r *Repo) List(ctx context.Context) ([]entity.Note, error) {
	var notes []entity.Note
	err := r.db.
		From(tableName).
		Select("id", "text", "created_at").
		Order(goqu.C("created_at").Asc()).
		Executor().
		ScanStructsContext(ctx, &notes)
	if err != nil {
		return nil, err
	}

	if notes == nil {
		return []entity.Note{}, nil
	}

	return notes, nil
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	result, err := r.db.
		Delete(tableName).
		Where(goqu.C("id").Eq(id)).
		Executor().
		ExecContext(ctx)
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
