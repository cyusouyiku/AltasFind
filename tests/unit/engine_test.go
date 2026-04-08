package unit_test

import (
	"context"
	"testing"
	"time"

	"atlasfind/internal/config"
	"atlasfind/internal/engine"
)

func TestEngineRunFindsMatches(t *testing.T) {
	tmp := t.TempDir()
	mustWriteFile(t, tmp+`/a.txt`, "first line\npanic happened\n")
	mustWriteFile(t, tmp+`/b.txt`, "all good here\n")
	mustWriteFile(t, tmp+`/c.log`, "panic again\n")

	cfg := config.Config{
		Root:        tmp,
		Pattern:     "panic",
		Workers:     2,
		Timeout:     2 * time.Second,
		Extensions:  []string{".txt", ".log"},
		MaxFileSize: 1024 * 1024,
	}

	eng, err := engine.New(cfg)
	if err != nil {
		t.Fatalf("engine.New returned error: %v", err)
	}

	results, err := eng.Run(context.Background())
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 files with matches, got %d", len(results))
	}
}
