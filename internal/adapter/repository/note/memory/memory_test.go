package memory_test

import (
	"testing"
	"time"

	"github.com/jian-hua-he/reference-app/internal/adapter/repository"
	"github.com/jian-hua-he/reference-app/internal/adapter/repository/note/memory"
	"github.com/jian-hua-he/reference-app/internal/entity"
	"github.com/jian-hua-he/reference-app/internal/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryRepo_List(t *testing.T) {
	testCases := map[string]struct {
		InitialNotes map[string]memory.Note

		Want    []entity.Note
		WantErr error
	}{
		"empty list": {
			InitialNotes: map[string]memory.Note{},
			Want:         []entity.Note{},
		},
		"list with notes": {
			InitialNotes: map[string]memory.Note{
				"1": {Text: "note 1", CreatedAt: time.Unix(10, 0).UTC()},
				"2": {Text: "note 2", CreatedAt: time.Unix(20, 0).UTC()},
			},
			Want: []entity.Note{
				{ID: "1", Text: "note 1", CreatedAt: time.Unix(10, 0).UTC()},
				{ID: "2", Text: "note 2", CreatedAt: time.Unix(20, 0).UTC()},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			repo := memory.NewRepoWithNotes(
				test.FakeUUID,
				test.FakeNow,
				tc.InitialNotes,
			)

			got, err := repo.List(t.Context())

			assert.ErrorIs(t, err, tc.WantErr)
			assert.ElementsMatch(t, got, tc.Want)
		})
	}
}

func TestMemoryRepo_Create(t *testing.T) {
	testCases := map[string]struct {
		Input string

		InitialNotes map[string]memory.Note

		Want         *entity.Note
		WantAllNotes []entity.Note
		WantErr      error
	}{
		"create note with text": {
			Input: "hello world",

			Want: &entity.Note{
				ID:        "fake-uuid-1234",
				Text:      "hello world",
				CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			WantAllNotes: []entity.Note{
				{
					ID:        "fake-uuid-1234",
					Text:      "hello world",
					CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				},
			},
		},
		"create note with spaces": {
			Input: "  spaced text  ",

			Want: &entity.Note{
				ID:        "fake-uuid-1234",
				Text:      "  spaced text  ",
				CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			WantAllNotes: []entity.Note{
				{
					ID:        "fake-uuid-1234",
					Text:      "  spaced text  ",
					CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				},
			},
		},
		"create empty note": {
			Input: "",

			InitialNotes: map[string]memory.Note{
				"1": {
					Text:      "Foo",
					CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				"2": {
					Text:      "Bar",
					CreatedAt: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},

			Want: &entity.Note{
				ID:        "fake-uuid-1234",
				Text:      "",
				CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			WantAllNotes: []entity.Note{
				{
					ID:        "1",
					Text:      "Foo",
					CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        "2",
					Text:      "Bar",
					CreatedAt: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        "fake-uuid-1234",
					Text:      "",
					CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			repo := memory.NewRepoWithNotes(
				test.FakeUUID,
				test.FakeNow,
				tc.InitialNotes,
			)

			got, err := repo.Create(t.Context(), tc.Input)
			gotList, listErr := repo.List(t.Context())
			require.NoError(t, listErr)

			assert.ErrorIs(t, err, tc.WantErr)
			assert.EqualValues(t, tc.Want, got)
			assert.ElementsMatch(t, tc.WantAllNotes, gotList)
		})
	}
}

func TestMemoryRepo_Delete(t *testing.T) {
	testCases := map[string]struct {
		InitialNotes map[string]memory.Note
		InputID      string

		Want    []entity.Note
		WantErr error
	}{
		"delete existing note": {
			InputID: "1",
			InitialNotes: map[string]memory.Note{
				"1": {
					Text:      "to delete",
					CreatedAt: time.Date(2024, 5, 5, 0, 0, 0, 0, time.UTC),
				},
				"2": {
					Text:      "to keep",
					CreatedAt: time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
				},
				"3": {
					Text:      "also to keep",
					CreatedAt: time.Date(2026, 6, 6, 0, 0, 0, 0, time.UTC),
				},
			},

			Want: []entity.Note{
				{
					ID:        "2",
					Text:      "to keep",
					CreatedAt: time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        "3",
					Text:      "also to keep",
					CreatedAt: time.Date(2026, 6, 6, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		"delete non-existent note": {
			InputID: "non-existent-id",
			InitialNotes: map[string]memory.Note{
				"1": {
					Text:      "to delete",
					CreatedAt: time.Date(2024, 5, 5, 0, 0, 0, 0, time.UTC),
				},
				"2": {
					Text:      "to keep",
					CreatedAt: time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
				},
				"3": {
					Text:      "also to keep",
					CreatedAt: time.Date(2026, 6, 6, 0, 0, 0, 0, time.UTC),
				},
			},

			Want: []entity.Note{
				{
					ID:        "1",
					Text:      "to delete",
					CreatedAt: time.Date(2024, 5, 5, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        "2",
					Text:      "to keep",
					CreatedAt: time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        "3",
					Text:      "also to keep",
					CreatedAt: time.Date(2026, 6, 6, 0, 0, 0, 0, time.UTC),
				},
			},
			WantErr: repository.ErrNotFound,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			repo := memory.NewRepoWithNotes(
				test.FakeUUID,
				test.FakeNow,
				tc.InitialNotes,
			)

			err := repo.Delete(t.Context(), tc.InputID)
			assert.ErrorIs(t, err, tc.WantErr)

			got, listErr := repo.List(t.Context())
			require.NoError(t, listErr)

			assert.ElementsMatch(t, tc.Want, got)
		})
	}
}
