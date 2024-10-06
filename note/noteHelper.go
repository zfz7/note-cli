package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func OpenNote(relativeWeek int, config NoteConfig) error {
	now := time.Now()
	//Find monday then add relative week
	daysToSubtract := ((int(now.Weekday()) + 6) % 7) + -relativeWeek*7
	monday := now.AddDate(0, 0, -daysToSubtract)
	noteName := monday.Format("2006-01-02") + "." + config.Extension
	notePath, _ := AppendHomeDirectory(config.Location + "/" + noteName)
	noteExists, err := FileExists(notePath)

	if err != nil {
		fmt.Println("Could not determine if file exists: " + notePath)
		return err
	}

	if !noteExists {
		err := CreateNewNote(notePath, err, config)
		if err != nil {
			fmt.Printf("Error creating new note")
		}
	}
	//Open the file
	cmd := exec.Command(config.Editor, notePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()

	if err != nil {
		fmt.Println("Error running command:", err)
		return err
	}
	return nil
}

func CreateNewNote(notePath string, err error, config NoteConfig) error {
	// Create the parent directories if they don't exist
	templateFile, err := ReadFile(config.Template)

	if err != nil {
		fmt.Printf("Error opening template file: %s\n", err)
		// Create empty note file
		err = WriteFile(notePath, []byte(""))
	}
	err = WriteFile(notePath, templateFile)
	if err != nil {
		fmt.Printf("Error creating new note at: "+notePath, err)
		return err
	}
	return nil
}
