package main

import (
	"context"
	"time"

	"github.com/jian-hua-he/ddd_notes/internal/entity"
)

type FakeApp struct{}

func NewFakeApp() *FakeApp {
	return &FakeApp{}
}

func (a *FakeApp) List(ctx context.Context) ([]entity.Note, error) {
	return []entity.Note{
		{
			ID:        "1",
			Text:      "This is a fake note 1",
			CreatedAt: time.Now(),
		},
		{
			ID:        "2",
			Text:      "This is a fake note 2",
			CreatedAt: time.Now(),
		},
	}, nil
}

func (a *FakeApp) Create(ctx context.Context, text string) (*entity.Note, error) {
	return &entity.Note{
		ID:        "fake_id",
		Text:      text,
		CreatedAt: time.Now(),
	}, nil
}

func (a *FakeApp) Delete(ctx context.Context, id string) error {
	return nil
}
