package main

import (
	"os"

	"github.com/0t4u/ragtag-backup/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	cmd.Execute()
}
