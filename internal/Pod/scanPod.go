package pod

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	schematics "github.com/DragonEmperor9480/yorasys/internal/Schematics"
)

func ScanAnamolies(reg schematics.Registry) {

	for _, valCache := range reg.Caches {
		cachePresent := false
		fmt.Printf("\nCache: %s (ID: %d)\n", valCache.Name, valCache.ID)

		for _, cachePath := range valCache.Paths {
			subPaths, err := handleFullPath(cachePath)
			if err != nil {
				fmt.Printf("Wrong Yaml data on %v, err: %v\n", cachePath, err)
			}

			if len(subPaths) == 0 {
				fmt.Printf("am never gonna execute but justtt lets see", cachePath)
			}
			exists, _, err := checkPath(cachePath)
			if err != nil {
				fmt.Printf("program.exe is meow meow %v\n", err)
				continue
			}
			if exists {
				cachePresent = true
				fmt.Printf("Found something\n", cachePath)
			} else {
				fmt.Printf("Meh didnt find a thing\n", cachePath)
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
