package walker

import (
	"context"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

// Options controls which files are emitted by the walker.
type Options struct {
	Root          string
	IncludeHidden bool
	Extensions    []string
	MaxFileSize   int64
}

// CollectFiles walks the filesystem and returns all matched files.
func CollectFiles(ctx context.Context, opts Options) ([]string, error) {
	files := make([]string, 0, 32)
	out := make(chan string, 64)
	errCh := make(chan error, 1)

	go func() {
		errCh <- Walk(ctx, opts, out)
		close(out)
	}()

	for path := range out {
		files = append(files, path)
	}

	sort.Strings(files)
	return files, <-errCh
}

// Walk streams file paths that match the provided options.
func Walk(ctx context.Context, opts Options, out chan<- string) error {
	if strings.TrimSpace(opts.Root) == "" {
		opts.Root = "."
	}

	extensions := make(map[string]struct{}, len(opts.Extensions))
	for _, ext := range opts.Extensions {
		extensions[strings.ToLower(ext)] = struct{}{}
	}

	return filepath.WalkDir(opts.Root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if path == opts.Root {
			return nil
		}

		name := entry.Name()
		if !opts.IncludeHidden && strings.HasPrefix(name, ".") {
			if entry.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if entry.IsDir() {
			return nil
		}

		if len(extensions) > 0 {
			ext := strings.ToLower(filepath.Ext(name))
			if _, ok := extensions[ext]; !ok {
				return nil
			}
		}

		if opts.MaxFileSize > 0 {
			info, err := entry.Info()
			if err == nil && info.Size() > opts.MaxFileSize {
				return nil
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- path:
			return nil
		}
	})
}
