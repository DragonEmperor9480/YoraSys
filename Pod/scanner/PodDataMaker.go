package pod

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	schematics "github.com/DragonEmperor9480/yorasys/Pod/Schematics"
)

const (
	archiveDir         = "archives"
	scanResultDirName  = "scanres"
	cleanResultDirName = "cleanres"
)

func WriteScanArchive(scanData schematics.ScanData) (string, error) {
	generatedAt := time.Now()
	scanResultDir := filepath.Join(archiveDir, scanResultDirName)
	cleanResultDir := filepath.Join(archiveDir, cleanResultDirName)
	if err := os.MkdirAll(cleanResultDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create clean result folder: %w", err)
	}

	archiveFolder := filepath.Join(scanResultDir, fmt.Sprintf("scan_%s", generatedAt.Format("20060102_150405")))
	if err := os.MkdirAll(archiveFolder, 0755); err != nil {
		return "", fmt.Errorf("failed to create archive folder: %w", err)
	}

	archive := schematics.ScanArchive{
		GeneratedAt:    generatedAt.Format(time.RFC3339),
		TotalFiles:     scanData.TotalFiles,
		TotalSizeBytes: scanData.TotalSizeBytes,
		ScannedPaths:   map[string]schematics.ScanArchivePath{},
	}

	rootPaths := make([]string, 0, len(scanData.ScannedPaths))
	for rootPath := range scanData.ScannedPaths {
		rootPaths = append(rootPaths, rootPath)
	}
	sort.Strings(rootPaths)

	for _, rootPath := range rootPaths {
		pathData := scanData.ScannedPaths[rootPath]
		files := make([]schematics.ScanArchiveEntry, 0, len(pathData.Files))
		for _, file := range pathData.Files {
			files = append(files, schematics.ScanArchiveEntry{
				Name:      file.Name,
				SizeBytes: file.SizeBytes,
			})
		}
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name < files[j].Name
		})

		archive.ScannedPaths[rootPath] = schematics.ScanArchivePath{
			TotalFiles:     pathData.TotalFiles,
			TotalSizeBytes: pathData.TotalSizeBytes,
			Files:          files,
		}
	}

	archivePath := filepath.Join(archiveFolder, "cache_files.json")

	file, err := os.Create(archivePath)
	if err != nil {
		return "", fmt.Errorf("failed to create archive json: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(archive); err != nil {
		return "", fmt.Errorf("failed to write archive json: %w", err)
	}

	return archivePath, nil
}
