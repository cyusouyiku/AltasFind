package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"atlasfind/internal/config"
	"atlasfind/internal/engine"
	"atlasfind/internal/result"
)

func main() {
	cfg := parseFlags()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	eng, err := engine.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		flag.Usage()
		os.Exit(2)
	}

	started := time.Now()
	results, err := eng.Run(ctx)
	if cfg.JSON {
		if encodeErr := json.NewEncoder(os.Stdout).Encode(results); encodeErr != nil {
			fmt.Fprintf(os.Stderr, "encode results: %v\n", encodeErr)
			os.Exit(1)
		}
	} else {
		printHuman(results)
		fmt.Fprintf(os.Stderr, "completed in %s\n", time.Since(started).Round(time.Millisecond))
	}

	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			fmt.Fprintln(os.Stderr, "search timed out")
		case errors.Is(err, context.Canceled):
			fmt.Fprintln(os.Stderr, "search cancelled")
		default:
			fmt.Fprintf(os.Stderr, "search failed: %v\n", err)
		}
		os.Exit(1)
	}
}

func parseFlags() config.Config {
	cfg := config.Config{}
	var extText string

	flag.StringVar(&cfg.Root, "root", ".", "root directory to search")
	flag.StringVar(&cfg.Pattern, "pattern", "", "regex or text pattern to search")
	flag.IntVar(&cfg.Workers, "workers", 0, "number of concurrent workers (default: CPU cores)")
	flag.BoolVar(&cfg.IgnoreCase, "i", false, "ignore case while matching")
	flag.BoolVar(&cfg.IgnoreCase, "ignore-case", false, "ignore case while matching")
	flag.BoolVar(&cfg.Literal, "literal", false, "treat pattern as literal text")
	flag.BoolVar(&cfg.IncludeHidden, "hidden", false, "include hidden files and folders")
	flag.DurationVar(&cfg.Timeout, "timeout", 30*time.Second, "overall timeout, e.g. 5s or 1m")
	flag.Int64Var(&cfg.MaxFileSize, "max-size", 2<<20, "skip files larger than this many bytes")
	flag.StringVar(&extText, "ext", "", "comma-separated file extensions to include, e.g. .go,.md")
	flag.BoolVar(&cfg.JSON, "json", false, "print JSON output")

	flag.Parse()
	cfg.Extensions = config.ParseExtensions(extText)
	return cfg
}

func printHuman(results []result.FileResult) {
	if len(results) == 0 {
		fmt.Println("no matches found")
		return
	}

	total := 0
	for _, file := range results {
		for _, match := range file.Matches {
			fmt.Printf("%s:%d: %s\n", file.Path, match.LineNumber, match.Line)
			total++
		}
	}

	fmt.Fprintf(os.Stderr, "found %d matches in %d files\n", total, len(results))
}
