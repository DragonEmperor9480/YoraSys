package pod

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const archiveDir = "archives"

type scanArchive struct {
	GeneratedAt    string                     `json:"generated_at"`
	TotalFiles     int                        `json:"total_files"`
	TotalSizeBytes int64                      `json:"total_size_bytes"`
	ScannedPaths   map[string]scanArchivePath `json:"scanned_paths"`
}

type scanArchivePath struct {
	TotalFiles     int                `json:"total_files"`
	TotalSizeBytes int64              `json:"total_size_bytes"`
	Files          []scanArchiveEntry `json:"files"`
}

type scanArchiveEntry struct {
	Name      string `json:"name"`
	SizeBytes int64  `json:"size_bytes"`
}

func WriteScanArchive(scanData ScanData) (string, error) {
	generatedAt := time.Now()
	archiveFolder := filepath.Join(archiveDir, fmt.Sprintf("scan_%s", generatedAt.Format("20060102_150405")))
	if err := os.MkdirAll(archiveFolder, 0755); err != nil {
		return "", fmt.Errorf("failed to create archive folder: %w", err)
	}

	archive := scanArchive{
		GeneratedAt:    generatedAt.Format(time.RFC3339),
		TotalFiles:     scanData.TotalFiles,
		TotalSizeBytes: scanData.TotalSizeBytes,
		ScannedPaths:   map[string]scanArchivePath{},
	}

	rootPaths := make([]string, 0, len(scanData.ScannedPaths))
	for rootPath := range scanData.ScannedPaths {
		rootPaths = append(rootPaths, rootPath)
	}
	sort.Strings(rootPaths)

	for _, rootPath := range rootPaths {
		pathData := scanData.ScannedPaths[rootPath]
		files := make([]scanArchiveEntry, 0, len(pathData.Files))
		for _, file := range pathData.Files {
			files = append(files, scanArchiveEntry{
				Name:      file.Name,
				SizeBytes: file.SizeBytes,
			})
		}
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name < files[j].Name
		})

		archive.ScannedPaths[rootPath] = scanArchivePath{
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
