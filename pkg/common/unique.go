package common

import (
	"fmt"
	"sort"
	"strings"
)

const maxUniqueLines = 5000

// UniqueCountLines counts duplicate lines and returns uniq -c style output.
func UniqueCountLines(content string) string {
	lines := strings.Split(content, "\n")
	counts := make(map[string]int, len(lines))
	for _, line := range lines {
		counts[line]++
	}

	type lineCount struct {
		line  string
		count int
	}
	pairs := make([]lineCount, 0, len(counts))
	for line, count := range counts {
		pairs = append(pairs, lineCount{line: line, count: count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].line < pairs[j].line
	})
	if len(pairs) > maxUniqueLines {
		pairs = pairs[:maxUniqueLines]
	}

	var b strings.Builder
	for _, p := range pairs {
		fmt.Fprintf(&b, "%7d %s\n", p.count, p.line)
	}
	return b.String()
}
