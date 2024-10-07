package main

import (
	"embed"
	"encoding/json"
	"fmt"
)

type ConfigHelper interface {
	Config() error
	ReadConfig() (NoteConfig, error)
}

type configHelper struct {
	fileHelper FileHelper
}

func NewConfigHelper(fileHelper FileHelper) ConfigHelper {
	return &configHelper{
		fileHelper: fileHelper,
	}
}

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

func (configHelper configHelper) Config() error {
	config, err := configHelper.ReadConfig()
	if err != nil { //Config file does not exist
		config, err = configHelper.writeDefaultConfig()
		if err != nil {
			return err
		}
	}

	//Open the file
	cleandConfigPath, err := configHelper.fileHelper.AppendHomeDirectory(ConfigPath)
	if err != nil {
		return err
	}

	err = configHelper.fileHelper.EditorOpenFile(config.Editor, cleandConfigPath)
	if err != nil {
		return err
	}
	return nil
}

func (configHelper configHelper) writeDefaultConfig() (NoteConfig, error) {
	bytes, _ := json.MarshalIndent(defaultConfig, "", "  ")
	err := configHelper.fileHelper.WriteFile(ConfigPath, bytes)
	if err != nil {
		fmt.Println("Error writing default config: "+string(bytes), err)
		return NoteConfig{}, err
	}

	staticTemplate, err := staticFiles.ReadFile("static/template.md")
	if err != nil {
		fmt.Println("Error reading default template:", err)
		return NoteConfig{}, err
	}
	err = configHelper.fileHelper.WriteFile(defaultConfig.Template, staticTemplate)
	if err != nil {
		fmt.Println("Error writing default template: "+string(bytes), err)
		return NoteConfig{}, err
	}
	fmt.Println("Default config written to " + ConfigPath + " to change config manually update config.json")
	return defaultConfig, nil
}

func (configHelper configHelper) ReadConfig() (NoteConfig, error) {
	config := NoteConfig{}
	data, err := configHelper.fileHelper.ReadFile(ConfigPath)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		defaultConfig, _ := json.Marshal(config)
		fmt.Println("Error Unmarshalling file at: "+ConfigPath+" expected format: "+string(defaultConfig), err)
		return config, err
	}
	return config, nil
}
