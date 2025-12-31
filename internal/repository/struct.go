package repository

import "time"

type RepoNote struct {
	ID        string
	Text      string
	CreatedAt time.Time
}
