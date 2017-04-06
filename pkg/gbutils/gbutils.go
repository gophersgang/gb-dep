package gbutils

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindInAncestorPath will go up from given path until it finds the required folder/file
func FindInAncestorPath(dir string, folderOrFile string) (string, error) {
	found := false
	foundPath := ""
	for {
		expectedPath := filepath.Join(dir, folderOrFile)
		if PathExist(expectedPath) {
			foundPath = expectedPath
			found = true
			break
		}
		next := filepath.Clean(filepath.Join(dir, ".."))
		if next == "/" {
			dir = "/"
			break
		}
		dir = next
	}
	if found {
		return foundPath, nil
	}
	return "", fmt.Errorf("%s not found in %s", folderOrFile, dir)
}

// PathExist is a quick way to check for folder/file existence
func PathExist(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
}

// IsFile checks whether give path is a file
func IsFile(path string) bool {
	if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
		return true
	}
	return false
}
