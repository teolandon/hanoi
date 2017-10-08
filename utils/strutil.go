package strutil

import (
	"fmt"
)

// WrapText wraps a long string into a slice of at most 'lines' strings, with
// the condition that small-ish words are preserved, and each string in
// the returned slice is at most 'limit' long.
func WrapText(text string, limit, lines int) []string {
	fmt.Println("Trying to wrap text at", limit, "length")
	ret := make([]string, lines)

	lineCount, i := 0, 0

	for i = limit; i < len(text) && lineCount < lines; i += limit {
		start := i - limit
		var cut bool
		fmt.Println("Index of last break char before", i, ":")
		i, cut = indexOfLastBreakChar(text, start, i+1)
		fmt.Println(i)
		fmt.Println()

		ret[lineCount] = text[start:i]
		lineCount++

		if cut {
			fmt.Println("cut")
			i++
		}
	}

	// Cover last case
	if lineCount < lines {
		ret[lineCount] = text[i-limit:]
	}

	return ret
}

func indexOfLastBreakChar(s string, start, end int) (i int, cut bool) {
	for i := end - 1; i >= start; i-- {
		switch s[i] {
		case ' ':
			return i, true
		case ',', '-', '.', '/', ')', '}', ':', ';', '"', '?', '>', '!':
			return i, false
		case '{', '(', '<':
			return i - 1, false
		}
	}
	return end - 1, false
}
