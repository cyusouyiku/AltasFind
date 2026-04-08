// Package searcher provides regular-expression based matching helpers.
package searcher

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"

	"atlasfind/internal/result"
)

// RegexSearcher searches content using a compiled regular expression.
type RegexSearcher struct {
	Re *regexp.Regexp
}

// Search returns all substring matches from a plain text blob.
func (r *RegexSearcher) Search(text string) ([]string, error) {
	if r == nil || r.Re == nil {
		return nil, nil
	}
	return r.Re.FindAllString(text, -1), nil
}

// FindMatches scans content line by line and reports matching lines.
func (r *RegexSearcher) FindMatches(content []byte) []result.Match {
	if r == nil || r.Re == nil {
		return nil
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	matches := make([]result.Match, 0)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		groups := r.Re.FindAllString(line, -1)
		if len(groups) == 0 {
			continue
		}

		matches = append(matches, result.Match{
			LineNumber: lineNumber,
			Line:       strings.TrimSpace(line),
			Groups:     groups,
		})
	}

	return matches
}
