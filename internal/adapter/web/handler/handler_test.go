package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jian-hua-he/reference-app/internal/adapter/web/handler"
	"github.com/jian-hua-he/reference-app/internal/adapter/web/router"
	"github.com/jian-hua-he/reference-app/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
)

func TestHandler_List(t *testing.T) {
	testCases := map[string]struct {
		MockAppListReturnNotes []entity.Note
		MockAppListReturnErr   error
		MockAppListCallTimes   int

		Want           handler.GetNotesResponse
		WantStatusCode int
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
			MockAppListReturnErr: nil,
			MockAppListCallTimes: 1,

			Want: handler.GetNotesResponse{
				Message: "successful",
				Payload: []handler.Note{
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
			},
			WantStatusCode: http.StatusOK,
		},
		"app error": {
			MockAppListReturnNotes: nil,
			MockAppListReturnErr:   assert.AnError,
			MockAppListCallTimes:   1,

			Want: handler.GetNotesResponse{
				Message: "failed to list notes",
			},
			WantStatusCode: http.StatusInternalServerError,
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
				Times(tc.MockAppListCallTimes)

			h := handler.NewHandler(app)

			e := echo.New()
			e.GET(router.UrlPathNote, h.GetNotes)

			req := httptest.NewRequest(http.MethodGet, router.UrlPathNote, nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.WantStatusCode, rec.Code)

			var got handler.GetNotesResponse
			err := json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)
			assert.Equal(t, tc.Want.Message, got.Message)
			assert.ElementsMatch(t, tc.Want.Payload, got.Payload)
		})
	}
}

func TestHandler_Create(t *testing.T) {
	testCases := map[string]struct {
		Input handler.PostNoteRequest

		MockAppCreateInput      string
		MockAppCreateReturnNote *entity.Note
		MockAppCreateReturnErr  error
		MockAppCreateCallTimes  int

		Want           handler.PostNoteResponse
		WantStatusCode int
	}{
		"successful create": {
			Input: handler.PostNoteRequest{
				Text: "new note",
			},

			MockAppCreateInput: "new note",
			MockAppCreateReturnNote: &entity.Note{
				ID:        "1",
				Text:      "new note",
				CreatedAt: time.Unix(10, 0).UTC(),
			},
			MockAppCreateCallTimes: 1,

			Want: handler.PostNoteResponse{
				Message: "successful",
				Payload: &handler.Note{
					ID:        "1",
					Text:      "new note",
					CreatedAt: time.Unix(10, 0).UTC(),
				},
			},
			WantStatusCode: http.StatusOK,
		},
		"app error": {
			Input: handler.PostNoteRequest{
				Text: "new note",
			},

			MockAppCreateInput:     "new note",
			MockAppCreateReturnErr: assert.AnError,
			MockAppCreateCallTimes: 1,

			Want: handler.PostNoteResponse{
				Message: "failed to create note",
			},
			WantStatusCode: http.StatusInternalServerError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			app := NewMockNoteApp(ctrl)
			app.EXPECT().
				Create(gomock.Any(), tc.MockAppCreateInput).
				Return(tc.MockAppCreateReturnNote, tc.MockAppCreateReturnErr).
				Times(tc.MockAppCreateCallTimes)

			h := handler.NewHandler(app)

			e := echo.New()
			e.POST(router.UrlPathNote, h.PostNote)

			body, err := json.Marshal(tc.Input)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, router.UrlPathNote, bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.WantStatusCode, rec.Code)

			var got handler.PostNoteResponse
			err = json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)
			assert.EqualValues(t, tc.Want, got)
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	testCases := map[string]struct {
		Input string

		MockAppDeleteInput     string
		MockAppDeleteReturnErr error
		MockAppDeleteCallTimes int

		Want           handler.DeleteNoteResponse
		WantStatusCode int
	}{
		"successful delete": {
			Input: "1",

			MockAppDeleteInput:     "1",
			MockAppDeleteCallTimes: 1,

			Want: handler.DeleteNoteResponse{
				Message: "successful",
			},
			WantStatusCode: http.StatusOK,
		},
		"app error": {
			Input: "2",

			MockAppDeleteInput:     "2",
			MockAppDeleteReturnErr: assert.AnError,
			MockAppDeleteCallTimes: 1,

			Want: handler.DeleteNoteResponse{
				Message: "failed to delete note",
			},
			WantStatusCode: http.StatusInternalServerError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			app := NewMockNoteApp(ctrl)
			app.EXPECT().
				Delete(gomock.Any(), tc.MockAppDeleteInput).
				Return(tc.MockAppDeleteReturnErr).
				Times(tc.MockAppDeleteCallTimes)

			h := handler.NewHandler(app)

			e := echo.New()
			e.DELETE(router.UrlPathNoteWithID, h.DeleteNote)

			path := fmt.Sprintf("%s/%s", router.UrlPathNote, tc.Input)
			req := httptest.NewRequest(http.MethodDelete, path, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.WantStatusCode, rec.Code)

			var got handler.DeleteNoteResponse
			err := json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)
			assert.EqualValues(t, tc.Want, got)
		})
	}
}
