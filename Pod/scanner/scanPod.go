package pod

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	schematics "github.com/DragonEmperor9480/yorasys/Pod/Schematics"
	"github.com/DragonEmperor9480/yorasys/Pod/service"
)

func ScanAnamolies(reg schematics.Registry) schematics.ScanData {
	scanData := schematics.ScanData{
		ScannedPaths: map[string]schematics.ScannedPathData{},
	}
	globalSeenFiles := map[string]bool{}

	for _, valCache := range reg.Caches {
		cachePresent := false
		var cacheTotalBytes int64
		var cacheMappedFiles int
		seen := map[string]bool{}
		seenFiles := map[string]bool{}
		fmt.Printf("\nCache: %s (ID: %d)\n", valCache.Name, valCache.ID)

		for _, cachePath := range valCache.Paths {
			expandedPath, missing := service.ExpandWindowsEnv(cachePath)
			if len(missing) > 0 {
				fmt.Printf("Unresolved env vars in %s: %v\n", cachePath, missing)
				continue
			}

			subPaths, err := service.HandleFullPath(expandedPath)
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

				exists, isDir, err := service.CheckPath(subPath)
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

				filePaths := make([]string, 0, len(pathFileSizes))
				for filePath := range pathFileSizes {
					filePaths = append(filePaths, filePath)
				}
				sort.Strings(filePaths)

				pathData := scanData.ScannedPaths[subPath]
				var pathTotalBytes int64
				var addedFiles int
				for _, filePath := range filePaths {
					sizeBytes := pathFileSizes[filePath]
					normalizedFilePath := strings.ToLower(filepath.Clean(filePath))
					if seenFiles[normalizedFilePath] {
						continue
					}
					seenFiles[normalizedFilePath] = true

					if globalSeenFiles[normalizedFilePath] {
						continue
					}
					globalSeenFiles[normalizedFilePath] = true

					pathData.Files = append(pathData.Files, schematics.ScannedFileData{
						Name:      service.RelativeFileName(subPath, filePath, isDir),
						SizeBytes: sizeBytes,
					})
					pathData.TotalFiles++
					pathData.TotalSizeBytes += sizeBytes

					scanData.TotalFiles++
					scanData.TotalSizeBytes += sizeBytes
					pathTotalBytes += sizeBytes
					addedFiles++
				}

				if pathData.TotalFiles > 0 {
					scanData.ScannedPaths[subPath] = pathData
				}
				cacheTotalBytes += pathTotalBytes
				cacheMappedFiles += addedFiles
				fmt.Printf("Found something: %s | size: %.2f MB (%d bytes) | files: %d\n", subPath, bytesToMB(pathTotalBytes), pathTotalBytes, addedFiles)
			}
		}
		fmt.Printf("hmmmmm %v | total_size: %.2f MB (%d bytes) | mapped_files: %d\n", cachePresent, bytesToMB(cacheTotalBytes), cacheTotalBytes, cacheMappedFiles)
	}

	return scanData
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
