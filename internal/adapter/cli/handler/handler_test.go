package handler_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/jian-hua-he/reference-app/internal/adapter/cli/handler"
	"github.com/jian-hua-he/reference-app/internal/entity"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_ListNotes(t *testing.T) {
	testCases := map[string]struct {
		MockAppListReturnNotes []entity.Note
		MockAppListReturnErr   error

		WantOutput string
		WantErr    error
	}{
		"successful list": {
			MockAppListReturnNotes: []entity.Note{
				{
					ID:        "1",
					Text:      "note 1",
					CreatedAt: time.Unix(10, 0).UTC(),
				},
				{
					ID:        "2",
					Text:      "note 2",
					CreatedAt: time.Unix(20, 0).UTC(),
				},
			},

			WantOutput: "ID: 1\nText: note 1\nCreated At: 1970-01-01 00:00:10 +0000 UTC\n\nID: 2\nText: note 2\nCreated At: 1970-01-01 00:00:20 +0000 UTC\n\n",
		},
		"empty list": {
			MockAppListReturnNotes: []entity.Note{},

			WantOutput: "No notes found.\n",
		},
		"app error": {
			MockAppListReturnErr: assert.AnError,

			WantErr: assert.AnError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			app := handler.NewMockNoteApp(ctrl)
			app.EXPECT().
				List(gomock.Any()).
				Return(tc.MockAppListReturnNotes, tc.MockAppListReturnErr).
				Times(1)

			var buf bytes.Buffer
			h := handler.NewHandler(app, &buf)

			err := h.ListNotes(t.Context())

			assert.ErrorIs(t, err, tc.WantErr)
			assert.Equal(t, tc.WantOutput, buf.String())
		})
	}
}

func TestHandler_CreateNote(t *testing.T) {
	testCases := map[string]struct {
		Input string

		MockAppCreateReturnNote *entity.Note
		MockAppCreateReturnErr  error

		WantOutput string
		WantErr    error
	}{
		"successful create": {
			Input: "new note",

			MockAppCreateReturnNote: &entity.Note{
				ID:        "1",
				Text:      "new note",
				CreatedAt: time.Unix(10, 0).UTC(),
			},

			WantOutput: "Note created successfully.\nID: 1\nText: new note\nCreated At: 1970-01-01 00:00:10 +0000 UTC\n",
		},
		"app error": {
			Input: "new note",

			MockAppCreateReturnErr: assert.AnError,

			WantErr: assert.AnError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			app := handler.NewMockNoteApp(ctrl)
			app.EXPECT().
				Create(gomock.Any(), tc.Input).
				Return(tc.MockAppCreateReturnNote, tc.MockAppCreateReturnErr).
				Times(1)

			var buf bytes.Buffer
			h := handler.NewHandler(app, &buf)

			err := h.CreateNote(t.Context(), tc.Input)

			assert.ErrorIs(t, err, tc.WantErr)
			assert.Equal(t, tc.WantOutput, buf.String())
		})
	}
}

func TestHandler_DeleteNote(t *testing.T) {
	testCases := map[string]struct {
		InputID string

		MockAppDeleteReturnErr error

		WantOutput string
		WantErr    error
	}{
		"successful delete": {
			InputID: "1",

			WantOutput: "Note deleted successfully.\n",
		},
		"app error": {
			InputID: "2",

			MockAppDeleteReturnErr: assert.AnError,

			WantErr: assert.AnError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			app := handler.NewMockNoteApp(ctrl)
			app.EXPECT().
				Delete(gomock.Any(), tc.InputID).
				Return(tc.MockAppDeleteReturnErr).
				Times(1)

			var buf bytes.Buffer
			h := handler.NewHandler(app, &buf)

			err := h.DeleteNote(t.Context(), tc.InputID)

			assert.ErrorIs(t, err, tc.WantErr)
			assert.Equal(t, tc.WantOutput, buf.String())
		})
	}
}
