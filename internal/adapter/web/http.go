package web

import "time"

type ErrResponse struct {
	Message string `json:"message"`
}

type ListNoteResponse struct {
	Message string `json:"message"`
	Payload []Note `json:"payload"`
}

type Note struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
