package gbutils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
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

// ComputeMD5 computes ... well.. md5 checksum of a file
func ComputeMD5(filePath string) (string, error) {
	// var result string
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	chksum := hex.EncodeToString(hash.Sum(nil))
	return chksum, nil
}

// ComputeMD5Content computes md5 checksum for a string
func ComputeMD5Content(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// ContainsStr checks for existance of a string in a slice
func ContainsStr(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
