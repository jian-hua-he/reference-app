package cli

import (
	"github.com/jian-hua-he/reference-app/internal/adapter/cli/handler"

	"github.com/spf13/cobra"
)

func NewRootCommand(h *handler.Handler) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "notes",
		Short: "Note management CLI",
	}

	rootCmd.AddCommand(
		newListCommand(h),
		newCreateCommand(h),
		newDeleteCommand(h),
	)

	return rootCmd
}

func newListCommand(h *handler.Handler) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all notes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.ListNotes(cmd.Context())
		},
	}
}

func newCreateCommand(h *handler.Handler) *cobra.Command {
	var text string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new note",
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.CreateNote(cmd.Context(), text)
		},
	}

	cmd.Flags().StringVar(&text, "text", "", "Note text")
	cmd.MarkFlagRequired("text")

	return cmd
}

func newDeleteCommand(h *handler.Handler) *cobra.Command {
	var id string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a note by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.DeleteNote(cmd.Context(), id)
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "Note ID")
	cmd.MarkFlagRequired("id")

	return cmd
}
