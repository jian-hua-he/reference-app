package handler

import "time"

type Note struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type (
	PostNoteRequest struct {
		Text string `json:"text"`
	}

	PostNoteResponse struct {
		Message string `json:"message"`
		Payload *Note  `json:"payload"`
	}
)

type GetNotesResponse struct {
	Message string `json:"message"`
	Payload []Note `json:"payload"`
}

type DeleteNoteResponse struct {
	Message string `json:"message"`
}
