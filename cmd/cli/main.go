package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jian-hua-he/reference-app/internal/application/note"
	"github.com/jian-hua-he/reference-app/internal/repository/note/memory"
	"github.com/jian-hua-he/reference-app/pkg/uuid"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	repo := memory.NewRepo(uuid.NewUUID, time.Now)
	app := note.NewNoteApp(repo)

	subcmd := os.Args[1]

	var err error
	switch subcmd {
	case "list":
		err = runList(app)
	case "create":
		err = runCreate(app, os.Args[2:])
	case "delete":
		err = runDelete(app, os.Args[2:])
	default:
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: notes <command>")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  list     List all notes")
	fmt.Fprintln(os.Stderr, "  create   Create a new note")
	fmt.Fprintln(os.Stderr, "  delete   Delete a note by ID")
}

func runList(app *note.NoteApp) error {
	ctx := context.Background()

	notes, err := app.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	if len(notes) == 0 {
		fmt.Println("No notes found.")
		return nil
	}

	for _, n := range notes {
		fmt.Printf("ID: %s\nText: %s\nCreated At: %s\n\n", n.ID, n.Text, n.CreatedAt)
	}

	return nil
}

func runCreate(app *note.NoteApp, args []string) error {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	text := fs.String("text", "", "Note text (required)")
	fs.Parse(args)

	if *text == "" {
		return fmt.Errorf("--text is required")
	}

	ctx := context.Background()

	n, err := app.Create(ctx, *text)
	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	fmt.Println("Note created successfully.")
	fmt.Printf("ID: %s\nText: %s\nCreated At: %s\n", n.ID, n.Text, n.CreatedAt)

	return nil
}

func runDelete(app *note.NoteApp, args []string) error {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	id := fs.String("id", "", "Note ID (required)")
	fs.Parse(args)

	if *id == "" {
		return fmt.Errorf("--id is required")
	}

	ctx := context.Background()

	err := app.Delete(ctx, *id)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	fmt.Println("Note deleted successfully.")

	return nil
}
