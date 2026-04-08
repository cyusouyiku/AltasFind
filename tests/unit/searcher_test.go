package unit_test

import (
	"testing"

	"atlasfind/internal/searcher"
)

func TestSearchWithRegex(t *testing.T) {
	matches, err := searcher.SearchWithRegex(`foo\d+`, "foo1 bar foo2 baz")
	if err != nil {
		t.Fatalf("SearchWithRegex returned error: %v", err)
	}

	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
}

func TestRegexSearcherFindMatches(t *testing.T) {
	rs, err := searcher.NewRegexSearcher(`(?i)atlas`)
	if err != nil {
		t.Fatalf("NewRegexSearcher returned error: %v", err)
	}

	results := rs.FindMatches([]byte("atlas\nAtlas project\nnone"))
	if len(results) != 2 {
		t.Fatalf("expected 2 line matches, got %d", len(results))
	}

	if results[0].LineNumber != 1 || results[1].LineNumber != 2 {
		t.Fatalf("unexpected line numbers: %+v", results)
	}
}
