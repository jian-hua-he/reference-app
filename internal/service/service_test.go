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
	testCases := map[string]struct {
		MockListReturnNotes []repository.RepoNote
		MockListReturnErr   error
		MockListCallTimes   int

		Want    []service.ServiceNote
		WantErr error
	}{
		"successful": {
			MockListReturnNotes: []repository.RepoNote{
				{ID: "1", Text: "a", CreatedAt: time.Unix(10, 0).UTC()},
				{ID: "2", Text: "b", CreatedAt: time.Unix(20, 0).UTC()},
			},
			MockListCallTimes: 1,

			Want: []service.ServiceNote{
				{ID: "1", Text: "a", CreatedAt: time.Unix(10, 0).UTC()},
				{ID: "2", Text: "b", CreatedAt: time.Unix(20, 0).UTC()},
			},
		},
		"repo error": {
			MockListReturnErr: assert.AnError,
			MockListCallTimes: 1,

			WantErr: assert.AnError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := service.NewMockNoteRepository(ctrl)
			svc := service.NewNoteService(repo)

			repo.EXPECT().
				List(gomock.Any()).
				Return(tc.MockListReturnNotes, tc.MockListReturnErr).
				Times(1)

			got, err := svc.List(t.Context())

			assert.ErrorIs(t, err, tc.WantErr)
			assert.Equal(t, tc.Want, got)
		})
	}
}

func TestNoteService_Create(t *testing.T) {
	testCases := map[string]struct {
		Input string

		MockCreateReturn    *repository.RepoNote
		MockCreateReturnErr error
		MockCreateCallTimes int

		Want    *service.ServiceNote
		WantErr error
	}{
		"passes text with space": {
			Input: "  hello  ",

			MockCreateReturn: &repository.RepoNote{
				ID:        "id-1",
				Text:      "  hello  ",
				CreatedAt: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
			},
			MockCreateCallTimes: 1,

			Want: &service.ServiceNote{
				ID:        "id-1",
				Text:      "  hello  ",
				CreatedAt: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
			},
		},
		"pass empty text": {
			Input: "",

			MockCreateReturn: &repository.RepoNote{
				ID:        "id-2",
				Text:      "",
				CreatedAt: time.Unix(20, 0).UTC(),
			},
			MockCreateCallTimes: 1,

			Want: &service.ServiceNote{
				ID:        "id-2",
				Text:      "",
				CreatedAt: time.Unix(20, 0).UTC(),
			},
		},
		"repo error": {
			Input: "x",

			MockCreateReturnErr: assert.AnError,
			MockCreateCallTimes: 1,

			WantErr: assert.AnError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := service.NewMockNoteRepository(ctrl)
			svc := service.NewNoteService(repo)

			repo.EXPECT().
				Create(gomock.Any(), tc.Input).
				Return(tc.MockCreateReturn, tc.MockCreateReturnErr).
				Times(1)

			got, err := svc.Create(t.Context(), tc.Input)

			assert.ErrorIs(t, err, tc.WantErr)
			assert.Equal(t, tc.Want, got)
		})
	}
}

func TestNoteService_Delete(t *testing.T) {
	testCases := map[string]struct {
		InputID string

		MockDeleteErr       error
		MockDeleteCallTimes int

		WantErr error
	}{
		"successful": {
			InputID: "1",

			MockDeleteCallTimes: 1,
		},
		"note not found": {
			InputID: "404",

			MockDeleteErr:       repository.ErrNotFound,
			MockDeleteCallTimes: 1,

			WantErr: service.ErrNotFound,
		},
		"repo error": {
			InputID: "x",

			MockDeleteErr:       assert.AnError,
			MockDeleteCallTimes: 1,

			WantErr: assert.AnError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := service.NewMockNoteRepository(ctrl)
			svc := service.NewNoteService(repo)

			repo.EXPECT().
				Delete(gomock.Any(), tc.InputID).
				Return(tc.MockDeleteErr).
				Times(1)

			err := svc.Delete(t.Context(), tc.InputID)
			assert.ErrorIs(t, err, tc.WantErr)
		})
	}
}
