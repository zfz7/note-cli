package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func WriteFile(path string, bytes []byte) error {
	cleanedPath, err := AppendHomeDirectory(path)
	if err != nil {
		fmt.Println("Error appending home directory to: "+path, err)
		return err
	}

	// Create the parent directories if they don't exist
	dir := filepath.Dir(cleanedPath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return err
	}

	// Create the config file
	file, err := os.Create(cleanedPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}
	return nil
}

func ReadFile(path string) ([]byte, error) {
	cleanedPath, err := AppendHomeDirectory(path)
	if err != nil {
		fmt.Println("Error appending home directory to: "+path, err)
		return nil, err
	}

	data, err := os.ReadFile(cleanedPath)
	if err != nil {
		fmt.Println("Error reading file at: "+cleanedPath, err)
		return nil, err
	}
	return data, nil
}

func FileExists(path string) (bool, error) {
	cleanedPath, err := AppendHomeDirectory(path)
	if err != nil {
		fmt.Println("Error appending home directory to: "+path, err)
		return false, err
	}
	_, err = os.Stat(cleanedPath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		fmt.Printf("Error checking file: "+cleanedPath, err)
		return false, err
	}
	return true, nil
}

func AppendHomeDirectory(filePath string) (string, error) {
	if strings.HasPrefix(filePath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error determining home directory:", err)
			return "", err
		}
		return filepath.Join(homeDir, filePath[2:]), nil
	}
	return filePath, nil
}
