package unit_test

import (
	"context"
	"path/filepath"
	"testing"

	"atlasfind/internal/walker"
)

func TestCollectFilesSkipsHiddenAndFiltersExtension(t *testing.T) {
	tmp := t.TempDir()
	visible := filepath.Join(tmp, "main.go")
	hiddenDir := filepath.Join(tmp, ".git")
	hiddenFile := filepath.Join(hiddenDir, "config")
	other := filepath.Join(tmp, "notes.txt")

	mustWriteFile(t, visible, "package main")
	mustWriteFile(t, hiddenFile, "ignored")
	mustWriteFile(t, other, "ignored")

	files, err := walker.CollectFiles(context.Background(), walker.Options{
		Root:          tmp,
		IncludeHidden: false,
		Extensions:    []string{".go"},
	})
	if err != nil {
		t.Fatalf("CollectFiles returned error: %v", err)
	}

	if len(files) != 1 || files[0] != visible {
		t.Fatalf("unexpected files: %#v", files)
	}
}
