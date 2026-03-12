package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jian-hua-he/reference-app/internal/application/note"
	"github.com/jian-hua-he/reference-app/internal/repository/note/memory"
	"github.com/jian-hua-he/reference-app/pkg/uuid"
)

func main() {
	repo := memory.NewRepo(uuid.NewUUID, time.Now)
	app := note.NewNoteApp(repo)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Note management CLI. Type 'help' for commands, 'exit' to quit.")
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		args := strings.Fields(scanner.Text())
		if len(args) == 0 {
			continue
		}

		cmd := args[0]

		switch cmd {
		case "exit", "quit":
			fmt.Println("Bye!")
			return
		case "help":
			printUsage()
		case "list":
			if err := runList(app); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
		case "create":
			text := strings.TrimSpace(strings.Join(args[1:], " "))
			if err := runCreate(app, text); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
		case "delete":
			if len(args) < 2 {
				fmt.Fprintln(os.Stderr, "Usage: delete <id>")
				continue
			}
			if err := runDelete(app, args[1]); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
			printUsage()
		}
	}
}

func printUsage() {
	fmt.Println("Commands:")
	fmt.Println("  list               List all notes")
	fmt.Println("  create <text>      Create a new note")
	fmt.Println("  delete <id>        Delete a note by ID")
	fmt.Println("  exit               Exit the CLI")
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

func runCreate(app *note.NoteApp, text string) error {
	if text == "" {
		return fmt.Errorf("text is required: create <text>")
	}

	ctx := context.Background()

	n, err := app.Create(ctx, text)
	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	fmt.Println("Note created successfully.")
	fmt.Printf("ID: %s\nText: %s\nCreated At: %s\n", n.ID, n.Text, n.CreatedAt)

	return nil
}

func runDelete(app *note.NoteApp, id string) error {
	ctx := context.Background()

	err := app.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	fmt.Println("Note deleted successfully.")

	return nil
}
