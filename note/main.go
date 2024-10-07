package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	fileHelper := NewFileHelper()
	configHelper := NewConfigHelper(fileHelper)
	noteHelper := NewNoteHelper(fileHelper)
	app := &cli.App{
		Name:    "note",
		Usage:   "Simple cli tool to help create notes from a template file and open notes from previous weeks.",
		Version: "v0.0.1",
		Commands: []*cli.Command{
			{
				Name:    "setup",
				Aliases: []string{"s"},
				Usage:   "Writes default config, if exists opens it",
				Action: func(cCtx *cli.Context) error {
					err := configHelper.Setup()
					if err != nil {
						fmt.Println("Error writing config")
						cli.Exit("Error writing config", 1)
					}
					return nil
				},
			},
			{
				Name:    "open",
				Aliases: []string{"o"},
				Usage:   "Open existing note or create new note from template",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "time",
						Usage:       "Index to open previous or next weeks notes",
						Aliases:     []string{"t"},
						DefaultText: "0",
					},
				},
				Action: func(cCtx *cli.Context) error {
					relativeWeek := cCtx.Int("time")
					config, err := configHelper.ReadConfig()
					if err != nil {
						fmt.Println("Missing config, please run 'note setup'")
						cli.Exit("Missing config, please run 'note setup'", 1)
						return err
					}
					err = noteHelper.OpenNote(relativeWeek, config)
					if err != nil {
						fmt.Println("Could not open note, check files exist or permissions")
						cli.Exit("Could not open note, check files exist or permissions", 1)
						return err
					}
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
