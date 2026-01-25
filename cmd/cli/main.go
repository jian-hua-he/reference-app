package main

import (
	"time"

	"github.com/jian-hua-he/reference-app/internal/adapter/cli"
	"github.com/jian-hua-he/reference-app/internal/application/note"
	"github.com/jian-hua-he/reference-app/internal/repository/note/memory"
	"github.com/jian-hua-he/reference-app/pkg/uuid"

	"github.com/rs/zerolog/log"
)

func main() {
	repo := memory.NewRepo(uuid.NewUUID, time.Now)
	app := note.NewNoteApp(repo)

	if err := cli.Run(app); err != nil {
		log.Error().Err(err).Msg("CLI exited with error")
	}
}
