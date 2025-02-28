package main

import (
	"fmt"
	"strings"
	"time"
)

var Now = time.Now

type NoteHelper interface {
	OpenNoteByInterval(relativeInterval int, config NoteConfig) error
	OpenNoteByFileName(file string, config NoteConfig) error
}

type noteHelper struct {
	fileHelper FileHelper
}

func NewNoteHelper(fileHelper FileHelper) NoteHelper {
	return &noteHelper{
		fileHelper: fileHelper,
	}
}

func (noteHelper noteHelper) OpenNoteByInterval(relativeInterval int, config NoteConfig) error {
	noteDate := noteHelper.getNoteDate(relativeInterval, config)
	noteName := noteDate.Format("2006-01-02") + "." + config.Extension
	return openNote(noteName, config, noteHelper)
}

func (noteHelper noteHelper) OpenNoteByFileName(file string, config NoteConfig) error {
	noteName := file + "." + config.Extension
	return openNote(noteName, config, noteHelper)
}

func openNote(noteName string, config NoteConfig, noteHelper noteHelper) error {
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

func (noteHelper noteHelper) getNoteDate(relativeInterval int, config NoteConfig) time.Time {
	now := Now()
	if strings.ToLower(config.Interval) == "day" {
		return now.AddDate(0, 0, relativeInterval)
	}
	if strings.ToLower(config.Interval) == "week" {
		//Find monday then add relative week
		daysToSubtract := ((int(now.Weekday()) + 6) % 7) + -relativeInterval*7
		monday := now.AddDate(0, 0, -daysToSubtract)
		return monday
	}
	if strings.ToLower(config.Interval) == "month" {
		//Find first of month then add relative month
		month := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return month.AddDate(0, relativeInterval, 0)
	}
	return now
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
