package gitio

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const gitCognitoFolder = ".GitCognito"

// GetHomeDir returns the path to the user's home directory
func GetHomeDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return homeDir, nil
}

// GetFilePath returns the absolute file path based on the relative path
func GetFilePath(relativePath string) (string, error) {
	homeDir, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(homeDir, gitCognitoFolder, relativePath)

	// Check if the path is outside the .GitCognito folder
	if !strings.HasPrefix(path, filepath.Join(homeDir, gitCognitoFolder)) {
		return "", fmt.Errorf("invalid file path: %s", relativePath)
	}

	return path, nil
}

// ReadFile reads a file and returns its contents as a string
func ReadFile(relativePath string) (string, error) {
	path, err := GetFilePath(relativePath)
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(bytes), nil
}

// CreateFile creates a file with the given contents
func CreateFile(relativePath, contents string) error {
	path, err := GetFilePath(relativePath)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	if err := ioutil.WriteFile(path, []byte(contents), 0644); err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	return nil
}

// DeleteFile deletes a file
func DeleteFile(relativePath string) error {
	path, err := GetFilePath(relativePath)
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
