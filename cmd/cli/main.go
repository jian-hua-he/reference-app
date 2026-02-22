package main

import (
	"os"
	"time"

	"github.com/jian-hua-he/reference-app/internal/adapter/cli"
	clihandler "github.com/jian-hua-he/reference-app/internal/adapter/cli/handler"
	"github.com/jian-hua-he/reference-app/internal/application/note"
	"github.com/jian-hua-he/reference-app/internal/repository/note/memory"
	"github.com/jian-hua-he/reference-app/pkg/uuid"

	"github.com/rs/zerolog/log"
)

func main() {
	repo := memory.NewRepo(uuid.NewUUID, time.Now)
	app := note.NewNoteApp(repo)
	h := clihandler.NewHandler(app, os.Stdout)

	rootCmd := cli.NewRootCommand(h)

	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("command failed")
		os.Exit(1)
	}
}
