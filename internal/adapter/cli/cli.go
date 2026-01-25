package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Run(app NoteApp) error {
	fmt.Println("Type 'help' for available commands, or 'exit' to quit")
	fmt.Println()

	rootCmd := newRootCmd(app)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		args := strings.Fields(input)
		rootCmd.SetArgs(args)
		if err := rootCmd.Execute(); err != nil {
			log.Error().Err(err).Msg("command execution failed")
		}
	}

	return nil
}

func newRootCmd(app NoteApp) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "note",
		Short:         "A tiny note CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Available commands:")
			fmt.Println("- list: List all notes")
			fmt.Println("- create --text=<text>: Create a new note")
			fmt.Println("- delete --id=<id>: Delete a note")
			fmt.Println("- help: Show this help message")
		},
	}

	rootCmd.AddCommand(newListCmd(app))
	rootCmd.AddCommand(newCreateCmd(app))
	rootCmd.AddCommand(newDeleteCmd(app))

	return rootCmd
}

func newListCmd(app NoteApp) *cobra.Command {
	return &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			notes, err := app.List(cmd.Context())
			if err != nil {
				return fmt.Errorf("list notes failed, err %w", err)
			}

			if len(notes) == 0 {
				fmt.Println("No notes found")
				return nil
			}

			for _, n := range notes {
				fmt.Printf("ID=%s, Text=%s\n", n.ID, n.Text)
			}

			return nil
		},
	}
}

func newCreateCmd(app NoteApp) *cobra.Command {
	var text string

	cmd := &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			if text == "" {
				return fmt.Errorf("--text is required")
			}

			n, err := app.Create(cmd.Context(), text)
			if err != nil {
				return fmt.Errorf("create note failed, err %w", err)
			}

			fmt.Printf("Note created: ID=%s, Text=%s\n", n.ID, n.Text)

			return nil
		},
	}

	cmd.Flags().StringVar(&text, "text", "", "Note text")
	return cmd
}

func newDeleteCmd(app NoteApp) *cobra.Command {
	var id string

	cmd := &cobra.Command{
		Use: "delete",
		RunE: func(cmd *cobra.Command, args []string) error {
			if id == "" {
				return fmt.Errorf("--id is required")
			}

			if err := app.Delete(cmd.Context(), id); err != nil {
				return fmt.Errorf("delete note failed, err %w", err)
			}

			fmt.Printf("Note deleted: ID=%s\n", id)

			return nil
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "Note ID")
	return cmd
}
