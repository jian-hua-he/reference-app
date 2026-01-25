package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewNoteCmd(app NoteApp) *cobra.Command {
	noteCmd := &cobra.Command{
		Use:           "note",
		Short:         "A tiny note CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("App launched, you can do following")
			fmt.Println("- note list")
			fmt.Println("- note create --text=\"Foo bar\"")
			fmt.Println("- note delete")
			fmt.Println("- note exit")
		},
	}

	noteCmd.AddCommand(NewCreateCmd(app))
	noteCmd.AddCommand(NewListCmd(app))
	noteCmd.AddCommand(NewDeleteCmd(app))
	noteCmd.AddCommand(NewExitCmd())

	return noteCmd
}

func NewCreateCmd(app NoteApp) *cobra.Command {
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
	if err := cmd.MarkFlagRequired("text"); err != nil {
		fmt.Printf("failed to mark text flag as required: %v\n", err)
	}

	return cmd
}

func NewDeleteCmd(app NoteApp) *cobra.Command {
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
	if err := cmd.MarkFlagRequired("id"); err != nil {
		fmt.Printf("failed to mark id flag as required: %v\n", err)
	}

	return cmd
}

func NewListCmd(app NoteApp) *cobra.Command {
	return &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			notes, err := app.List(cmd.Context())
			if err != nil {
				return fmt.Errorf("list notes failed, err %w", err)
			}

			for _, n := range notes {
				fmt.Printf("ID=%s, Text=%s\n", n.ID, n.Text)
			}

			return nil
		},
	}
}

func NewExitCmd() *cobra.Command {
	return &cobra.Command{
		Use: "exit",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Exiting the application. Goodbye!")
			os.Exit(1)
		},
	}
}
