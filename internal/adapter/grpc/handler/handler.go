package handler

import (
	"context"
	"errors"

	notev1 "github.com/jian-hua-he/reference-app/proto/note/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	notev1.UnimplementedNoteServiceServer
	noteApp NoteApp
}

func NewHandler(noteApp NoteApp) *Handler {
	return &Handler{
		noteApp: noteApp,
	}
}

func (h *Handler) ListNotes(ctx context.Context, req *notev1.ListNotesRequest) (*notev1.ListNotesResponse, error) {
	notes, err := h.noteApp.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list notes: %v", err)
	}

	pbNotes := make([]*notev1.Note, len(notes))
	for i, n := range notes {
		pbNotes[i] = &notev1.Note{
			Id:        n.ID,
			Text:      n.Text,
			CreatedAt: timestamppb.New(n.CreatedAt),
		}
	}

	return &notev1.ListNotesResponse{Notes: pbNotes}, nil
}

func (h *Handler) CreateNote(ctx context.Context, req *notev1.CreateNoteRequest) (*notev1.CreateNoteResponse, error) {
	note, err := h.noteApp.Create(ctx, req.GetText())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create note: %v", err)
	}

	return &notev1.CreateNoteResponse{
		Note: &notev1.Note{
			Id:        note.ID,
			Text:      note.Text,
			CreatedAt: timestamppb.New(note.CreatedAt),
		},
	}, nil
}

func (h *Handler) DeleteNote(ctx context.Context, req *notev1.DeleteNoteRequest) (*notev1.DeleteNoteResponse, error) {
	if err := h.noteApp.Delete(ctx, req.GetId()); err != nil {
		return nil, errors.Join(err, status.Error(codes.Internal, "failed to delete note"))
	}

	return &notev1.DeleteNoteResponse{}, nil
}
