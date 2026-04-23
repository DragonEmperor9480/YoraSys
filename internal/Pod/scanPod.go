package pod

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	schematics "github.com/DragonEmperor9480/yorasys/internal/Schematics"
)

func ScanAnamolies(reg schematics.Registry) {
	for _, valCache := range reg.Caches {
		cachePresent := false
		var cacheTotalBytes int64
		seen := map[string]bool{}
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
				continue
			}

			for _, subPath := range subPaths {
				normalizedPath := strings.ToLower(filepath.Clean(subPath))
				if seen[normalizedPath] {
					continue
				}
				seen[normalizedPath] = true

				exists, isDir, err := checkPath(subPath)
				if err != nil {
					fmt.Printf("program.exe is meow meow %v\n", err)
					continue
				}
				if !exists {
					fmt.Printf("Meh didnt find a thing: %s\n", subPath)
					continue
				}

				cachePresent = true
				sizeBytes, err := pathSize(subPath, isDir)
				if err != nil {
					fmt.Printf("Found something: %s | size error: %v\n", subPath, err)
					continue
				}

				cacheTotalBytes += sizeBytes
				fmt.Printf("Found something: %s | size: %.2f MB (%d bytes)\n", subPath, bytesToMB(sizeBytes), sizeBytes)
			}
		}
		fmt.Printf("hmmmmm %v | total_size: %.2f MB (%d bytes)\n", cachePresent, bytesToMB(cacheTotalBytes), cacheTotalBytes)
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

func folderSize(path string) (int64, error) {
	var total int64

	err := filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip inaccessible entries
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}
		total += info.Size()
		return nil
	})

	if err != nil {
		return 0, err
	}
	return total, nil
}

func pathSize(path string, isDir bool) (int64, error) {
	if isDir {
		return folderSize(path)
	}

	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func bytesToMB(bytes int64) float64 {
	return float64(bytes) / (1024 * 1024)
}
