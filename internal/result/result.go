package result

// Match represents one matched line in a file.
type Match struct {
	LineNumber int      `json:"line_number"`
	Line       string   `json:"line"`
	Groups     []string `json:"groups,omitempty"`
}

// FileResult groups all matches for a file.
type FileResult struct {
	Path    string  `json:"path"`
	Matches []Match `json:"matches"`
}

// Count returns the number of matches in the file result.
func (r FileResult) Count() int {
	return len(r.Matches)
}
