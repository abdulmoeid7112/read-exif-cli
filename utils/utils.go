package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	Dir  = "dir"
	File = "file"
)

func IsPathExists(path string) (string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New(fmt.Sprintf("Path or file '%s' does not exist\n", path))
		} else {
			return "", err
		}
	}

	if fileInfo.IsDir() {
		return Dir, nil
	}

	return File, nil
}

func IsImage(path string) bool {
	extension := filepath.Ext(path)
	imageExtensions := []string{".jpg", ".jpeg", ".png"}
	for _, ext := range imageExtensions {
		if ext == extension {
			return true
		}
	}

	fmt.Println(fmt.Sprintf("Not a valid image file format: %v", path))
	return false
}
