package main

import (
	"errors"
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
		Version: "v0.0.2",
		Commands: []*cli.Command{
			{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Writes default config, if exists opens it",
				Action: func(cCtx *cli.Context) error {
					err := configHelper.Config()
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
						Name:        "interval",
						Usage:       "Relative interval to open previous or next interval's notes",
						Aliases:     []string{"i"},
						DefaultText: "0",
					},
					&cli.StringFlag{
						Name:    "file",
						Usage:   "Open note by file name",
						Aliases: []string{"f"},
					},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("interval") == "" && cCtx.String("file") == "" {
						fmt.Println("flags --interval and --file cannot both be set")
						cli.Exit("flags --interval and --file cannot both be set", 1)
						return errors.New("flags --interval and --file cannot both be set")
					}

					relativeInterval := cCtx.Int("interval")
					file := cCtx.String("file")
					config, err := configHelper.ReadConfig()
					if err != nil {
						fmt.Println("Missing config, please run 'note config'")
						cli.Exit("Missing config, please run 'note config'", 1)
						return err
					}
					if file != "" {
						err = noteHelper.OpenNoteByFileName(file, config)
					} else {
						err = noteHelper.OpenNoteByInterval(relativeInterval, config)
					}
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
