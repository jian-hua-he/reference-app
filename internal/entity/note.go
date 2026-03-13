package entity

import "time"

type Note struct {
	ID        string    `db:"id"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
}
