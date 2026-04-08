package result

import (
	"sort"
	"sync"
)

// Collector stores file results in a concurrency-safe way.
type Collector struct {
	mu      sync.Mutex
	results []FileResult
}

func NewCollector() *Collector {
	return &Collector{results: make([]FileResult, 0, 16)}
}

func (c *Collector) Add(item FileResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.results = append(c.results, item)
}

func (c *Collector) All() []FileResult {
	c.mu.Lock()
	defer c.mu.Unlock()

	items := make([]FileResult, len(c.results))
	copy(items, c.results)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Path < items[j].Path
	})
	return items
}

func (c *Collector) TotalMatches() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	total := 0
	for _, item := range c.results {
		total += len(item.Matches)
	}
	return total
}
