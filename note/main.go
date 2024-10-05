package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//go:embed static/*
var staticFiles embed.FS

type Config struct {
	Editor        string `json:"editor"`
	Template      string `json:"template"`
	FileExtension string `json:"file_extension"`
	NotesFolder   string `json:"notes_folder"`
}

const ConfigPath = "~/.config/note/config.json"
const TemplatePath = "~/.config/note/template.md"

func main() {
	if len(os.Args) != 1 && len(os.Args) != 2 {
		fmt.Println("Invalid number of arguments, accept 0 or 1")
		os.Exit(-1)
	}

	if len(os.Args) == 2 && os.Args[1] == "setup" {
		flushDefaultConfig()
		os.Exit(0)
	}

	relativeWeek := getRelativeWeek()
	config := getConfig()
	openNote(relativeWeek, config)
}

func openNote(relativeWeek int, config Config) {
	now := time.Now()
	//Find monday then add relative week
	daysToSubtract := ((int(now.Weekday()) + 6) % 7) + -relativeWeek*7
	monday := now.AddDate(0, 0, -daysToSubtract)
	noteName := monday.Format("2006-01-02") + "." + config.FileExtension
	notePath, err := convertHomeDir(config.NotesFolder + "/" + noteName)
	if err != nil {
		fmt.Println("Invalid path: " + config.NotesFolder + "/" + noteName)
		os.Exit(-1)
	}
	// Check if the file exists
	_, err = os.Stat(notePath)
	if os.IsNotExist(err) {
		err := CreateNewNote(notePath, err, config)
		if err != nil {
			fmt.Printf("Error creating new note")
		}
	} else if err != nil {
		fmt.Printf("Error checking file: %s\n", err)
	}
	//Open the file
	cmd := exec.Command(config.Editor, notePath)
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}
	os.Exit(1)
}

func CreateNewNote(notePath string, err error, config Config) error {
	// Create the parent directories if they don't exist
	dir := filepath.Dir(notePath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return err
	}
	templatePath, err := convertHomeDir(config.Template)
	if err != nil {
		fmt.Println("Error reading template path:", err)
		return err
	}
	templateFile, err := os.Open(templatePath)

	if err != nil {
		fmt.Printf("Error opening template file: %s\n", err)
		// Create empty note file
		file, err := os.Create(notePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return err
		}
		defer file.Close()
		_, err = file.WriteString("")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return err
		}
		return nil
	}

	defer templateFile.Close()

	noteFile, err := os.Create(notePath)
	if err != nil {
		fmt.Printf("Error creating new note file: %s\n", err)
		return err
	}
	defer noteFile.Close()

	_, err = io.Copy(noteFile, templateFile)
	if err != nil {
		fmt.Printf("Error copying file: %s\n", err)
		return err
	}
	return nil
}

func flushDefaultConfig() {
	configPath, err := convertHomeDir(ConfigPath)
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	// Create the parent directories if they don't exist
	dir := filepath.Dir(configPath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return
	}

	// Create the config file
	file, err := os.Create(configPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	defaultConfig := Config{
		Editor:        "vim",
		Template:      TemplatePath,
		FileExtension: "md",
		NotesFolder:   "~/notes",
	}
	bytes, _ := json.Marshal(defaultConfig)
	_, err = file.Write(bytes)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}

	staticTemplate, err := staticFiles.ReadFile("static/template.md")
	if err != nil {
		fmt.Println("Error reading default template:", err)
	}
	templatePath, err := convertHomeDir(TemplatePath)
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}
	template, err := os.Create(templatePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	_, err = template.Write(staticTemplate)
	if err != nil {
		fmt.Println("Error writing template:", err)
	}
	defer template.Close()
}

func getConfig() Config {
	config := Config{
		Editor:        "vim",
		Template:      TemplatePath,
		FileExtension: "md",
		NotesFolder:   "~/notes",
	}
	configPath, err := convertHomeDir(ConfigPath)
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return config
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error reading config file at: "+configPath, err)
		return config
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		defaultConfig, _ := json.Marshal(config)
		fmt.Println("Error Unmarshalling file at: "+configPath+" expected format: "+string(defaultConfig), err)
		return config
	}
	return config
}

func getRelativeWeek() int {
	relativeWeek := 0
	if len(os.Args) == 2 {
		i, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not convert to int"+os.Args[1])
			os.Exit(-1)
		}
		relativeWeek = i
	}
	return relativeWeek
}

func convertHomeDir(filePath string) (string, error) {
	if strings.HasPrefix(filePath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			return "", err
		}
		return filepath.Join(homeDir, filePath[2:]), nil
	}
	return filePath, nil
}
