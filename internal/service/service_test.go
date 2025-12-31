package service_test

import (
	"testing"
	"time"

	"github.com/jian-hua-he/ddd_notes/internal/repository"
	"github.com/jian-hua-he/ddd_notes/internal/service"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNoteService_List(t *testing.T) {
	t.Parallel()

	type want struct {
		errIs error
		notes []service.ServiceNote
	}
	cases := map[string]struct {
		repoNotes []repository.RepoNote
		repoErr   error
		want      want
	}{
		"ok maps repo notes to service notes": {
			repoNotes: []repository.RepoNote{
				{ID: "1", Text: "a", CreatedAt: time.Unix(10, 0).UTC()},
				{ID: "2", Text: "b", CreatedAt: time.Unix(20, 0).UTC()},
			},
			want: want{
				errIs: nil,
				notes: []service.ServiceNote{
					{ID: "1", Text: "a", CreatedAt: time.Unix(10, 0).UTC()},
					{ID: "2", Text: "b", CreatedAt: time.Unix(20, 0).UTC()},
				},
			},
		},
		"repo error bubbles up": {
			repoErr: assert.AnError,
			want: want{
				errIs: assert.AnError,
				notes: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := service.NewMockNoteRepository(ctrl)
			svc := service.NewNoteService(repo)

			repo.EXPECT().
				List(gomock.Any()).
				Return(tc.repoNotes, tc.repoErr).
				Times(1)

			got, err := svc.List(t.Context())

			assert.ErrorIs(t, err, tc.want.errIs)
			assert.Equal(t, tc.want.notes, got)
		})
	}
}

func TestNoteService_Create(t *testing.T) {
	t.Parallel()

	type want struct {
		errIs error
		note  *service.ServiceNote
	}

	cases := map[string]struct {
		inputText string

		repoReturn *repository.RepoNote
		repoErr    error

		want want
	}{
		"ok passes text through as-is": {
			inputText: "  hello  ",
			repoReturn: &repository.RepoNote{
				ID:        "id-1",
				Text:      "  hello  ",
				CreatedAt: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
			},
			want: want{
				errIs: nil,
				note: &service.ServiceNote{
					ID:        "id-1",
					Text:      "  hello  ",
					CreatedAt: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		"ok allows empty text (no validation in MVP)": {
			inputText: "",
			repoReturn: &repository.RepoNote{
				ID:        "id-2",
				Text:      "",
				CreatedAt: time.Unix(20, 0).UTC(),
			},
			want: want{
				errIs: nil,
				note: &service.ServiceNote{
					ID:        "id-2",
					Text:      "",
					CreatedAt: time.Unix(20, 0).UTC(),
				},
			},
		},
		"repo error bubbles up": {
			inputText: "x",
			repoErr:   assert.AnError,
			want: want{
				errIs: assert.AnError,
				note:  nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := service.NewMockNoteRepository(ctrl)
			svc := service.NewNoteService(repo)

			repo.EXPECT().
				Create(gomock.Any(), tc.inputText).
				Return(tc.repoReturn, tc.repoErr).
				Times(1)

			got, err := svc.Create(t.Context(), tc.inputText)

			assert.ErrorIs(t, err, tc.want.errIs)
			assert.Equal(t, tc.want.note, got)
		})
	}
}

func TestNoteService_Delete(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		inputID string

		repoErr error
		wantErr error
	}{
		"ok": {
			inputID: "1",
			repoErr: nil,
			wantErr: nil,
		},
		"maps repo not found to service not found": {
			inputID: "404",
			repoErr: repository.ErrNotFound,
			wantErr: service.ErrNotFound,
		},
		"repo error bubbles up": {
			inputID: "x",
			repoErr: assert.AnError,
			wantErr: assert.AnError,
		},
		"ok allows empty id (no validation in MVP)": {
			inputID: "",
			repoErr: nil,
			wantErr: nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := service.NewMockNoteRepository(ctrl)
			svc := service.NewNoteService(repo)

			repo.EXPECT().
				Delete(gomock.Any(), tc.inputID).
				Return(tc.repoErr).
				Times(1)

			err := svc.Delete(t.Context(), tc.inputID)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
