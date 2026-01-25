package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/jian-hua-he/ddd_notes/internal/adapter/cli"
	"github.com/jian-hua-he/ddd_notes/internal/application/note"
	"github.com/jian-hua-he/ddd_notes/internal/repository/note/memory"
	"github.com/jian-hua-he/ddd_notes/pkg/uuid"
)

func main() {
	repo := memory.NewRepo(uuid.NewUUID, time.Now)
	app := note.NewNoteApp(repo)

	fmt.Println("Type 'help' for available commands, or 'exit' to quit")
	fmt.Println()

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

		if input == "exit" || input == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		// Parse input into arguments
		args := strings.Fields(input)

		// Create a new command for each execution
		cmd := cli.NewNoteCmd(app)
		cmd.SetArgs(args)

		if err := cmd.Execute(); err != nil {
			log.Error().Err(err).Msg("command execution failed")
		}

		fmt.Println()
	}
}
