package gitio

import (
	"io/ioutil"
	"os"
)

// ReadFile reads a file and returns its contents as a string
func ReadFile(path string) (string, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CreateFile creates a file with the given contents
func CreateFile(path string, contents string) error {
	err := ioutil.WriteFile(path, []byte(contents), 0644)
	return err
}

// DeleteFile deletes a file
func DeleteFile(path string) error {
	err := os.Remove(path)
	return err
}
