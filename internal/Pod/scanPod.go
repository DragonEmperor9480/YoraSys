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

func ScanAnamolies(reg schematics.Registry) map[string]int64 {
	fileSizeMap := map[string]int64{}

	for _, valCache := range reg.Caches {
		cachePresent := false
		var cacheTotalBytes int64
		seen := map[string]bool{}
		seenFiles := map[string]bool{}
		cacheFileSizeMap := map[string]int64{}
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
				pathFileSizes, err := collectPathSizes(subPath, isDir)
				if err != nil {
					fmt.Printf("Found something: %s | size error: %v\n", subPath, err)
					continue
				}

				var pathTotalBytes int64
				var addedFiles int
				for filePath, sizeBytes := range pathFileSizes {
					normalizedFilePath := strings.ToLower(filepath.Clean(filePath))
					if seenFiles[normalizedFilePath] {
						continue
					}
					seenFiles[normalizedFilePath] = true

					cacheFileSizeMap[filePath] = sizeBytes
					fileSizeMap[filePath] = sizeBytes
					pathTotalBytes += sizeBytes
					addedFiles++
				}

				cacheTotalBytes += pathTotalBytes
				fmt.Printf("Found something: %s | size: %.2f MB (%d bytes) | files: %d\n", subPath, bytesToMB(pathTotalBytes), pathTotalBytes, addedFiles)
			}
		}
		fmt.Printf("hmmmmm %v | total_size: %.2f MB (%d bytes) | mapped_files: %d\n", cachePresent, bytesToMB(cacheTotalBytes), cacheTotalBytes, len(cacheFileSizeMap))
	}

	return fileSizeMap
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

func collectPathSizes(path string, isDir bool) (map[string]int64, error) {
	pathSizes := map[string]int64{}

	if !isDir {
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		pathSizes[path] = info.Size()
		return pathSizes, nil
	}

	err := filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
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
		pathSizes[filePath] = info.Size()
		return nil
	})

	if err != nil {
		return nil, err
	}
	return pathSizes, nil
}

func bytesToMB(bytes int64) float64 {
	return float64(bytes) / (1024 * 1024)
}
