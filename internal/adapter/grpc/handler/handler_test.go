package handler_test

import (
	"testing"
	"time"

	"github.com/jian-hua-he/reference-app/internal/adapter/grpc/handler"
	"github.com/jian-hua-he/reference-app/internal/entity"
	notev1 "github.com/jian-hua-he/reference-app/proto/note/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestHandler_ListNotes(t *testing.T) {
	testCases := map[string]struct {
		MockAppListReturnNotes []entity.Note
		MockAppListReturnErr   error

		Want     *notev1.ListNotesResponse
		WantCode codes.Code
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

			Want: &notev1.ListNotesResponse{
				Notes: []*notev1.Note{
					{
						Id:        "1",
						Text:      "note 1",
						CreatedAt: timestamppb.New(time.Unix(10, 0).UTC()),
					},
					{
						Id:        "2",
						Text:      "note 2",
						CreatedAt: timestamppb.New(time.Unix(20, 0).UTC()),
					},
				},
			},
			WantCode: codes.OK,
		},
		"empty list": {
			MockAppListReturnNotes: []entity.Note{},

			Want: &notev1.ListNotesResponse{
				Notes: []*notev1.Note{},
			},
			WantCode: codes.OK,
		},
		"app error": {
			MockAppListReturnErr: assert.AnError,

			WantCode: codes.Internal,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			app := NewMockNoteApp(ctrl)
			app.EXPECT().
				List(gomock.Any()).
				Return(tc.MockAppListReturnNotes, tc.MockAppListReturnErr).
				Times(1)

			h := handler.NewHandler(app)

			got, err := h.ListNotes(t.Context(), &notev1.ListNotesRequest{})

			if tc.WantCode != codes.OK {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.WantCode, st.Code())
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Len(t, got.Notes, len(tc.Want.Notes))
			for i, want := range tc.Want.Notes {
				assert.Equal(t, want.Id, got.Notes[i].Id)
				assert.Equal(t, want.Text, got.Notes[i].Text)
				assert.Equal(t, want.CreatedAt.AsTime(), got.Notes[i].CreatedAt.AsTime())
			}
		})
	}
}

func TestHandler_CreateNote(t *testing.T) {
	testCases := map[string]struct {
		Input string

		MockAppCreateReturnNote *entity.Note
		MockAppCreateReturnErr  error

		Want     *notev1.CreateNoteResponse
		WantCode codes.Code
	}{
		"successful create": {
			Input: "new note",

			MockAppCreateReturnNote: &entity.Note{
				ID:        "1",
				Text:      "new note",
				CreatedAt: time.Unix(10, 0).UTC(),
			},

			Want: &notev1.CreateNoteResponse{
				Note: &notev1.Note{
					Id:        "1",
					Text:      "new note",
					CreatedAt: timestamppb.New(time.Unix(10, 0).UTC()),
				},
			},
			WantCode: codes.OK,
		},
		"app error": {
			Input: "new note",

			MockAppCreateReturnErr: assert.AnError,

			WantCode: codes.Internal,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			app := NewMockNoteApp(ctrl)
			app.EXPECT().
				Create(gomock.Any(), tc.Input).
				Return(tc.MockAppCreateReturnNote, tc.MockAppCreateReturnErr).
				Times(1)

			h := handler.NewHandler(app)

			got, err := h.CreateNote(t.Context(), &notev1.CreateNoteRequest{Text: tc.Input})

			if tc.WantCode != codes.OK {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.WantCode, st.Code())
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.Want.Note.Id, got.Note.Id)
			assert.Equal(t, tc.Want.Note.Text, got.Note.Text)
			assert.Equal(t, tc.Want.Note.CreatedAt.AsTime(), got.Note.CreatedAt.AsTime())
		})
	}
}

func TestHandler_DeleteNote(t *testing.T) {
	testCases := map[string]struct {
		InputID string

		MockAppDeleteReturnErr error

		WantCode codes.Code
		WantErr  error
	}{
		"successful delete": {
			InputID: "1",

			WantCode: codes.OK,
		},
		"app error": {
			InputID: "2",

			MockAppDeleteReturnErr: assert.AnError,

			WantCode: codes.Internal,
			WantErr:  assert.AnError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			app := NewMockNoteApp(ctrl)
			app.EXPECT().
				Delete(gomock.Any(), tc.InputID).
				Return(tc.MockAppDeleteReturnErr).
				Times(1)

			h := handler.NewHandler(app)

			_, err := h.DeleteNote(t.Context(), &notev1.DeleteNoteRequest{Id: tc.InputID})

			// if tc.WantCode != codes.OK {
			// 	require.Error(t, err)
			// 	st, ok := status.FromError(err)
			// 	require.True(t, ok)
			// 	assert.Equal(t, tc.WantCode, st.Code())
			// 	return
			// }
			st, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, tc.WantCode, st.Code())

			assert.ErrorIs(t, err, tc.WantErr)
		})
	}
}
