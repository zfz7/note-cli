package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

//go:embed static/*
var staticFiles embed.FS

type NoteConfig struct {
	Editor    string `json:"editor"`
	Location  string `json:"location"`
	Template  string `json:"template"`
	Extension string `json:"extension"`
}

const ConfigPath = "~/.config/note/config.json"
const DefaultTemplate = "~/.config/note/template.md"
const DefaultLocation = "~/notes"
const DefaultEditor = "vim"
const DefaultExtension = "md"

var defaultConfig = NoteConfig{
	Editor:    DefaultEditor,
	Location:  DefaultLocation,
	Template:  DefaultTemplate,
	Extension: DefaultExtension,
}

func WriteDefaultConfig() error {
	bytes, _ := json.Marshal(defaultConfig)
	err := WriteFile(ConfigPath, bytes)
	if err != nil {
		fmt.Println("Error writing default config: "+string(bytes), err)
		return err
	}

	staticTemplate, err := staticFiles.ReadFile("static/template.md")
	if err != nil {
		fmt.Println("Error reading default template:", err)
		return err
	}
	err = WriteFile(defaultConfig.Template, staticTemplate)
	if err != nil {
		fmt.Println("Error writing default template: "+string(bytes), err)
		return err
	}
	return nil
}

func ReadConfig() (NoteConfig, error) {
	config := NoteConfig{}
	data, err := ReadFile(ConfigPath)
	err = json.Unmarshal(data, &config)
	if err != nil {
		defaultConfig, _ := json.Marshal(config)
		fmt.Println("Error Unmarshalling file at: "+ConfigPath+" expected format: "+string(defaultConfig), err)
		return config, err
	}
	return config, nil
}

func ReadRelativeWeek() (int, error) {
	relativeWeek := 0
	if len(os.Args) == 2 {
		i, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not convert: "+os.Args[1]+" to int")
			return 0, err
		}
		relativeWeek = i
	}
	return relativeWeek, nil
}
