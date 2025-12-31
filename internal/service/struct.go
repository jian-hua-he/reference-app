package service

import "time"

type ServiceNote struct {
	ID        string
	Text      string
	CreatedAt time.Time
}
