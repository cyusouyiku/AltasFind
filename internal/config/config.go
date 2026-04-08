package config

import (
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	apperrors "atlasfind/pkg/errors"
)

// Config stores all runtime options for a search session.
type Config struct {
	Root          string
	Pattern       string
	Workers       int
	IgnoreCase    bool
	Literal       bool
	IncludeHidden bool
	Timeout       time.Duration
	MaxFileSize   int64
	Extensions    []string
	JSON          bool
}

// Normalize fills sensible defaults when options are omitted.
func (c *Config) Normalize() {
	if strings.TrimSpace(c.Root) == "" {
		c.Root = "."
	}
	if c.Workers <= 0 {
		c.Workers = runtime.NumCPU()
	}
	if c.Timeout <= 0 {
		c.Timeout = 30 * time.Second
	}
	if c.MaxFileSize <= 0 {
		c.MaxFileSize = 2 << 20 // 2 MiB
	}
}

// Validate checks if the configuration is usable.
func (c Config) Validate() error {
	if strings.TrimSpace(c.Pattern) == "" {
		return apperrors.ErrPatternRequired
	}

	info, err := os.Stat(c.Root)
	if err != nil || !info.IsDir() {
		return apperrors.ErrRootNotFound
	}

	if c.Workers <= 0 {
		return apperrors.ErrInvalidWorkers
	}

	return nil
}

// EffectivePattern returns the actual regexp pattern to compile.
func (c Config) EffectivePattern() string {
	pattern := c.Pattern
	if c.Literal {
		pattern = regexp.QuoteMeta(pattern)
	}
	if c.IgnoreCase {
		pattern = "(?i)" + pattern
	}
	return pattern
}

// ParseExtensions converts a comma-separated list into normalized extensions.
func ParseExtensions(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	seen := make(map[string]struct{}, len(parts))
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		ext := strings.ToLower(strings.TrimSpace(part))
		if ext == "" {
			continue
		}
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		if _, exists := seen[ext]; exists {
			continue
		}
		seen[ext] = struct{}{}
		result = append(result, ext)
	}

	return result
}
