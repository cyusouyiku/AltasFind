package engine

import (
	"bytes"
	"context"
	"errors"
	"os"
	"sync"

	"atlasfind/internal/config"
	"atlasfind/internal/result"
	"atlasfind/internal/searcher"
	"atlasfind/internal/walker"
	"atlasfind/pkg/logger"
)

// Engine coordinates walking, file reading, and content matching.
type Engine struct {
	cfg      config.Config
	searcher searcher.Searcher
}

// New creates a configured search engine.
func New(cfg config.Config) (*Engine, error) {
	cfg.Normalize()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	rs, err := searcher.NewRegexSearcher(cfg.EffectivePattern())
	if err != nil {
		return nil, err
	}

	return &Engine{cfg: cfg, searcher: rs}, nil
}

// Run executes the search pipeline and returns all matched files.
func (e *Engine) Run(ctx context.Context) ([]result.FileResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if e.cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.cfg.Timeout)
		defer cancel()
	}

	collector := result.NewCollector()
	pathCh := make(chan string, e.cfg.Workers*2)
	resultCh := make(chan result.FileResult, e.cfg.Workers)
	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	for i := 0; i < e.cfg.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.consume(ctx, pathCh, resultCh)
		}()
	}

	go func() {
		defer close(pathCh)
		err := walker.Walk(ctx, walker.Options{
			Root:          e.cfg.Root,
			IncludeHidden: e.cfg.IncludeHidden,
			Extensions:    e.cfg.Extensions,
			MaxFileSize:   e.cfg.MaxFileSize,
		}, pathCh)
		if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			select {
			case errCh <- err:
			default:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for item := range resultCh {
		collector.Add(item)
	}

	select {
	case err := <-errCh:
		return collector.All(), err
	default:
	}

	if err := ctx.Err(); err != nil {
		return collector.All(), err
	}

	return collector.All(), nil
}

func (e *Engine) consume(ctx context.Context, paths <-chan string, out chan<- result.FileResult) {
	for {
		select {
		case <-ctx.Done():
			return
		case path, ok := <-paths:
			if !ok {
				return
			}

			content, err := os.ReadFile(path)
			if err != nil {
				logger.Debugf("skip unreadable file %s: %v", path, err)
				continue
			}
			if bytes.IndexByte(content, 0) >= 0 {
				continue
			}

			matches := e.searcher.FindMatches(content)
			if len(matches) == 0 {
				continue
			}

			select {
			case <-ctx.Done():
				return
			case out <- result.FileResult{Path: path, Matches: matches}:
			}
		}
	}
}
