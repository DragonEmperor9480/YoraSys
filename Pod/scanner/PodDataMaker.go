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
	GeneratedAt    string             `json:"generated_at"`
	TotalFiles     int                `json:"total_files"`
	TotalSizeBytes int64              `json:"total_size_bytes"`
	Files          []scanArchiveEntry `json:"files"`
}

type scanArchiveEntry struct {
	Path      string `json:"path"`
	SizeBytes int64  `json:"size_bytes"`
}

func WriteScanArchive(fileSizeMap map[string]int64) (string, error) {
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create archive folder: %w", err)
	}

	paths := make([]string, 0, len(fileSizeMap))
	for path := range fileSizeMap {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	generatedAt := time.Now()
	archive := scanArchive{
		GeneratedAt: generatedAt.Format(time.RFC3339),
		TotalFiles:  len(paths),
		Files:       make([]scanArchiveEntry, 0, len(paths)),
	}

	for _, path := range paths {
		sizeBytes := fileSizeMap[path]
		archive.TotalSizeBytes += sizeBytes
		archive.Files = append(archive.Files, scanArchiveEntry{
			Path:      path,
			SizeBytes: sizeBytes,
		})
	}

	fileName := fmt.Sprintf("scan_%s.json", generatedAt.Format("20060102_150405"))
	archivePath := filepath.Join(archiveDir, fileName)

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
