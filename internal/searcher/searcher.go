package searcher

import (
	"regexp"

	"atlasfind/internal/result"
)

// Searcher finds matches from a file's content.
type Searcher interface {
	FindMatches(content []byte) []result.Match
}

// NewRegexSearcher compiles and returns a RegexSearcher.
func NewRegexSearcher(pattern string) (*RegexSearcher, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &RegexSearcher{Re: re}, nil
}

// SearchWithRegex is a convenience helper for raw string matching.
func SearchWithRegex(pattern string, text string) ([]string, error) {
	searcher, err := NewRegexSearcher(pattern)
	if err != nil {
		return nil, err
	}
	return searcher.Search(text)
}
