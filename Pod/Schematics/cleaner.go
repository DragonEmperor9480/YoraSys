package schematics

type ScanData struct {
	TotalFiles     int
	TotalSizeBytes int64
	ScannedPaths   map[string]ScannedPathData
}

type ScannedPathData struct {
	TotalFiles     int
	TotalSizeBytes int64
	Files          []ScannedFileData
}

type ScannedFileData struct {
	Name      string
	SizeBytes int64
}

type ScanArchive struct {
	GeneratedAt    string                     `json:"generated_at"`
	TotalFiles     int                        `json:"total_files"`
	TotalSizeBytes int64                      `json:"total_size_bytes"`
	ScannedPaths   map[string]ScanArchivePath `json:"scanned_paths"`
}

type ScanArchivePath struct {
	TotalFiles     int                `json:"total_files"`
	TotalSizeBytes int64              `json:"total_size_bytes"`
	Files          []ScanArchiveEntry `json:"files"`
}

type ScanArchiveEntry struct {
	Name      string `json:"name"`
	SizeBytes int64  `json:"size_bytes"`
}

type CleanSelection struct {
	ScanArchivePath string              `json:"scan_archive_path"`
	CreatedAt       string              `json:"created_at,omitempty"`
	DryRun          bool                `json:"dry_run,omitempty"`
	SelectedPaths   map[string][]string `json:"selected_paths"`
}

type CleanResult struct {
	CleanedAt              string            `json:"cleaned_at"`
	DryRun                 bool              `json:"dry_run"`
	ScanArchivePath        string            `json:"scan_archive_path"`
	SelectionPath          string            `json:"selection_path"`
	TotalSelectedFiles     int               `json:"total_selected_files"`
	TotalSelectedSizeBytes int64             `json:"total_selected_size_bytes"`
	WouldDeleteFiles       int               `json:"would_delete_files"`
	DeletedFiles           int               `json:"deleted_files"`
	FailedFiles            int               `json:"failed_files"`
	SkippedFiles           int               `json:"skipped_files"`
	TotalDeletedBytes      int64             `json:"total_deleted_bytes"`
	Results                []CleanFileResult `json:"results"`
}

type CleanFileResult struct {
	RootPath  string `json:"root_path"`
	Name      string `json:"name"`
	FullPath  string `json:"full_path"`
	SizeBytes int64  `json:"size_bytes"`
	Status    string `json:"status"`
	Error     string `json:"error,omitempty"`
}
