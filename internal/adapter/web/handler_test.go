package web_test

import (
	"encoding/json"
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

		Want           web.ListNoteResponse
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

			Want: web.ListNoteResponse{
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
			e.GET(web.NoteListRoute, handler.ListNotes)

			req := httptest.NewRequest(http.MethodGet, web.NoteListRoute, nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.WantStatusCode, rec.Code)

			var got web.ListNoteResponse
			err := json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)
			assert.ElementsMatch(t, tc.Want.Payload, got.Payload)
		})
	}
}
