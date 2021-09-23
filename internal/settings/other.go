package settings

// LibraryExplorer holds settings for the library explorer view.
type LibraryExplorer struct {
	DividerPosition float32  `json:"divider_position"`
	OpenRowKeys     []string `json:"open_row_keys"`
}

// FileRef holds a path to a file and an offset for all page references within that file.
type FileRef struct {
	Path   string `json:"path"`
	Offset int    `json:"offset"`
}

// ExportInfo holds information about a recent export so that it can be redone quickly.
type ExportInfo struct {
	TemplatePath string `json:"template_path"`
	ExportPath   string `json:"export_path"`
	LastUsed     int64  `json:"last_used"`
}
