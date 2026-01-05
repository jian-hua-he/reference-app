package repository

import "github.com/google/uuid"

func NewUUiD() string {
	return uuid.New().String()
}
