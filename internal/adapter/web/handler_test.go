package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jian-hua-he/ddd_notes/internal/adapter/web"
	"github.com/jian-hua-he/ddd_notes/internal/entity"
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

		Want           web.GetNotesResponse
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

			Want: web.GetNotesResponse{
				Message: "successful",
				Payload: []web.Note{
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

			Want: web.GetNotesResponse{
				Message: "failed to list notes",
			},
			WantStatusCode: http.StatusInternalServerError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			app := web.NewMockNoteApp(ctrl)
			app.EXPECT().
				List(gomock.Any()).
				Return(tc.MockAppListReturnNotes, tc.MockAppListReturnErr).
				Times(tc.MockAppListCallTimes)

			handler := web.NewHandler(app)

			e := echo.New()
			e.GET(web.UrlPathNote, handler.GetNotes)

			req := httptest.NewRequest(http.MethodGet, web.UrlPathNote, nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.WantStatusCode, rec.Code)

			var got web.GetNotesResponse
			err := json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)
			assert.Equal(t, tc.Want.Message, got.Message)
			assert.ElementsMatch(t, tc.Want.Payload, got.Payload)
		})
	}
}

func TestHandler_Create(t *testing.T) {
	testCases := map[string]struct {
		Input web.PostNoteRequest

		MockAppCreateInput      string
		MockAppCreateReturnNote *entity.Note
		MockAppCreateReturnErr  error
		MockAppCreateCallTimes  int

		Want           web.PostNoteResponse
		WantStatusCode int
	}{
		"successful create": {
			Input: web.PostNoteRequest{
				Text: "new note",
			},

			MockAppCreateInput: "new note",
			MockAppCreateReturnNote: &entity.Note{
				ID:        "1",
				Text:      "new note",
				CreatedAt: time.Unix(10, 0).UTC(),
			},
			MockAppCreateCallTimes: 1,

			Want: web.PostNoteResponse{
				Message: "successful",
				Payload: &web.Note{
					ID:        "1",
					Text:      "new note",
					CreatedAt: time.Unix(10, 0).UTC(),
				},
			},
			WantStatusCode: http.StatusOK,
		},
		"app error": {
			Input: web.PostNoteRequest{
				Text: "new note",
			},

			MockAppCreateInput:     "new note",
			MockAppCreateReturnErr: assert.AnError,
			MockAppCreateCallTimes: 1,

			Want: web.PostNoteResponse{
				Message: "failed to create note",
			},
			WantStatusCode: http.StatusInternalServerError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			app := web.NewMockNoteApp(ctrl)
			app.EXPECT().
				Create(gomock.Any(), tc.MockAppCreateInput).
				Return(tc.MockAppCreateReturnNote, tc.MockAppCreateReturnErr).
				Times(tc.MockAppCreateCallTimes)

			handler := web.NewHandler(app)

			e := echo.New()
			e.POST(web.UrlPathNote, handler.PostNote)

			body, err := json.Marshal(tc.Input)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, web.UrlPathNote, bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.WantStatusCode, rec.Code)

			var got web.PostNoteResponse
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

		Want           web.DeleteNoteResponse
		WantStatusCode int
	}{
		"successful delete": {
			Input: "1",

			MockAppDeleteInput:     "1",
			MockAppDeleteCallTimes: 1,

			Want: web.DeleteNoteResponse{
				Message: "successful",
			},
			WantStatusCode: http.StatusOK,
		},
		"app error": {
			Input: "2",

			MockAppDeleteInput:     "2",
			MockAppDeleteReturnErr: assert.AnError,
			MockAppDeleteCallTimes: 1,

			Want: web.DeleteNoteResponse{
				Message: "failed to delete note",
			},
			WantStatusCode: http.StatusInternalServerError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			app := web.NewMockNoteApp(ctrl)
			app.EXPECT().
				Delete(gomock.Any(), tc.MockAppDeleteInput).
				Return(tc.MockAppDeleteReturnErr).
				Times(tc.MockAppDeleteCallTimes)

			handler := web.NewHandler(app)

			e := echo.New()
			e.DELETE(web.UrlPathNoteWithID, handler.DeleteNote)

			path := fmt.Sprintf("%s/%s", web.UrlPathNote, tc.Input)
			req := httptest.NewRequest(http.MethodPost, path, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.WantStatusCode, rec.Code)

			var got web.DeleteNoteResponse
			err := json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)
			assert.EqualValues(t, tc.Want, got)
		})
	}
}
