package main

import (
	"fmt"
	"time"
)

var Now = time.Now

type NoteHelper interface {
	OpenNote(relativeWeek int, config NoteConfig) error
}

type noteHelper struct {
	fileHelper FileHelper
}

func NewNoteHelper(fileHelper FileHelper) NoteHelper {
	return &noteHelper{
		fileHelper: fileHelper,
	}
}

func (noteHelper noteHelper) OpenNote(relativeWeek int, config NoteConfig) error {
	now := Now()
	//Find monday then add relative week
	daysToSubtract := ((int(now.Weekday()) + 6) % 7) + -relativeWeek*7
	monday := now.AddDate(0, 0, -daysToSubtract)
	noteName := monday.Format("2006-01-02") + "." + config.Extension
	notePath, _ := noteHelper.fileHelper.AppendHomeDirectory(config.Location + "/" + noteName)
	noteExists, err := noteHelper.fileHelper.FileExists(notePath)

	if err != nil {
		fmt.Println("Could not determine if file exists: " + notePath)
		return err
	}

	if !noteExists {
		err := noteHelper.createNewNote(notePath, err, config)
		if err != nil {
			fmt.Printf("Error creating new note")
		}
	}
	err = noteHelper.fileHelper.EditorOpenFile(config.Editor, notePath)
	if err != nil {
		return err
	}
	return nil
}

func (noteHelper noteHelper) createNewNote(notePath string, err error, config NoteConfig) error {
	// Create the parent directories if they don't exist
	templateFile, err := noteHelper.fileHelper.ReadFile(config.Template)

	if err != nil {
		fmt.Printf("Error opening template file: %s\n", err)
		// Create empty note file
		err = noteHelper.fileHelper.WriteFile(notePath, []byte(""))
	}
	err = noteHelper.fileHelper.WriteFile(notePath, templateFile)
	if err != nil {
		fmt.Printf("Error creating new note at: "+notePath, err)
		return err
	}
	return nil
}
