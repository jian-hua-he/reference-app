package entity

import "time"

type Note struct {
	ID        string
	Text      string
	CreatedAt time.Time
}
