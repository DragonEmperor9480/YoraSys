package service

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func CheckPath(path string) (exists bool, isDir bool, err error) {
	info, err := os.Stat(path)
	if err == nil {
		return true, info.IsDir(), nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, false, nil
	}
	return false, false, err
}

//HandleFullPath is used to handle the paths which Contains Asterik(*) Symbol
func HandleFullPath(path string) ([]string, error) {
	if !strings.Contains(path, "*") {
		return []string{path}, nil
	}

	found, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}

	return found, nil
}

var winEnv = regexp.MustCompile(`%([A-Za-z0-9_]+)%`)

//ExpandWindowsEnv handles windows envinronment variable paths like %TEMP%,%APPDATA% etc
func ExpandWindowsEnv(path string) (string, []string) {
	unresolved := []string{}

	expanded := winEnv.ReplaceAllStringFunc(path, func(s string) string {
		key := strings.Trim(s, "%")
		val := os.Getenv(key)
		if val == "" {
			unresolved = append(unresolved, key)
			return s
		}
		return val
	})
	return expanded, unresolved
}

//RelativeFileName converts a full file path into the name/path we store inside JSON.
func RelativeFileName(rootPath string, filePath string, isDir bool) string {
	if !isDir {
		return filepath.Base(filePath)
	}

	relPath, err := filepath.Rel(rootPath, filePath)
	if err != nil || relPath == "." || strings.HasPrefix(relPath, "..") {
		return filePath
	}
	return relPath
}
