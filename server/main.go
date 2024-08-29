package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli"

	"github.com/chat-merger/merger/server/internal"
)

func main() {
	app := &internal.App{}
	cliApp := &cli.App{
		Name:         filepath.Base(os.Args[0]),
		HelpName:     filepath.Base(os.Args[0]),
		Usage:        "A new cli application",
		Writer:       os.Stdout,
		BashComplete: cli.DefaultAppComplete,
		Action:       app.Start,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "config",
				Usage:       "config file",
				Value:       "config.toml",
				Destination: &app.ConfigPath,
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
