package strutil

// WrapText wraps a long string into a slice of at most 'lines' strings, with
// the condition that small-ish words are preserved, and each string in
// the returned slice is at most 'limit' long.
func WrapText(text string, limit, lines int) []string {
	ret := make([]string, lines)

	lineCount, i := 0, 0

	for i = limit; i < len(text) && lineCount < lines; i += limit {
		start := i - limit
		i, _ = indexOfLastBreakChar(text, start, i+1)
		ret[lineCount] = text[start:i]
		lineCount++
	}

	// Cover last case
	if lineCount < lines {
		ret[lineCount] = text[i-limit:]
	}

	return ret
}

func indexOfLastBreakChar(s string, start, end int) (i int, cut bool) {
	for i := end; i >= start; i-- {
		if s[i] == ' ' {
			return i, false
		}
	}
	return end - 1, true
}
