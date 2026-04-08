package benchmark_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"atlasfind/internal/config"
	"atlasfind/internal/engine"
)

func BenchmarkEngineRun(b *testing.B) {
	tmp := b.TempDir()
	for i := 0; i < 200; i++ {
		content := "all good\n"
		if i%10 == 0 {
			content = "panic detected in worker\n"
		}
		path := filepath.Join(tmp, fmt.Sprintf("file-%03d.log", i))
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			b.Fatalf("failed to write benchmark file: %v", err)
		}
	}

	cfg := config.Config{
		Root:        tmp,
		Pattern:     "panic",
		Workers:     4,
		Timeout:     5 * time.Second,
		Extensions:  []string{".log"},
		MaxFileSize: 1024 * 1024,
	}

	eng, err := engine.New(cfg)
	if err != nil {
		b.Fatalf("engine.New returned error: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := eng.Run(context.Background()); err != nil {
			b.Fatalf("Run returned error: %v", err)
		}
	}
}
