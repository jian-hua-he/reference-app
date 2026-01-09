package main

import (
	"context"

	"github.com/jian-hua-he/ddd_notes/internal/entity"
)

type FakeApp struct{}

func NewFakeApp() *FakeApp {
	return &FakeApp{}
}

func (a *FakeApp) List(ctx context.Context) ([]entity.Note, error) {
	return nil, nil
}

func (a *FakeApp) Create(ctx context.Context, text string) (*entity.Note, error) {
	return nil, nil
}

func (a *FakeApp) Delete(ctx context.Context, id string) error {
	return nil
}
