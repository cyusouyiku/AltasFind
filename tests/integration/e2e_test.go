package integration_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIEndToEnd(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, "sample.txt")
	if err := os.WriteFile(file, []byte("hello\nneedle here\nbye\n"), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	cmd := exec.Command("go", "run", "./cmd/gofind", "-root", tmp, "-pattern", "needle")
	cmd.Dir = filepath.Join("..", "..")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\n%s", err, string(out))
	}

	output := string(out)
	if !strings.Contains(output, "needle here") {
		t.Fatalf("expected CLI output to contain match, got: %s", output)
	}
}
