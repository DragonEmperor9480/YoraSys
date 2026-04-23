package pod

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	schematics "github.com/DragonEmperor9480/yorasys/internal/Schematics"
)

func ScanAnamolies(reg schematics.Registry) {

	for _, valCache := range reg.Caches {
		cachePresent := false
		fmt.Printf("\nCache: %s (ID: %d)\n", valCache.Name, valCache.ID)

		for _, cachePath := range valCache.Paths {

			expandedPath, missing := expandWindowsEnv(cachePath)
			if len(missing) > 0 {
				fmt.Printf("Unresolved env vars in %s: %v\n", cachePath, missing)
				continue
			}

			subPaths, err := handleFullPath(expandedPath)
			if err != nil {
				fmt.Printf("Wrong Yaml data on %v, err: %v\n", cachePath, err)
				continue
			}

			if len(subPaths) == 0 {
				fmt.Printf("am never gonna execute but justtt lets see... path: %v", cachePath)
			}

			for _, subPath := range subPaths {
				exists, _, err := checkPath(subPath)
				if err != nil {
					fmt.Printf("program.exe is meow meow %v\n", err)
					continue
				}
				if exists {
					cachePresent = true
					fmt.Printf("Found something: %s\n", subPath)
				} else {
					fmt.Printf("Meh didnt find a thing: %s\n", subPath)
				}
			}
		}
		fmt.Printf("hmmmmm %v\n", cachePresent)
	}
}

func checkPath(path string) (exists bool, isDir bool, err error) {
	info, err := os.Stat(path)
	if err == nil {
		return true, info.IsDir(), nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, false, nil
	}
	return false, false, err
}

func handleFullPath(path string) ([]string, error) {
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

func expandWindowsEnv(path string) (string, []string) {
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
