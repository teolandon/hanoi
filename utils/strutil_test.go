package utils

import "testing"

func TestWrapText(t *testing.T) {
	cases := []struct {
		s            string
		limit, lines int
		want         []string
	}{
		{"This is me, doing stuff right now", 15, 6, make([]string, 6)},
	}
	for _, c := range cases {
		got := WrapText(c.s, c.limit, c.lines)
		for _, str := range got {
			t.Errorf("%v", str)
		}
	}
}

func TestIndexOfLastBreakChar(t *testing.T) {
	cases := []struct {
		s          string
		start, end int
		want       int
	}{
		{"This is me, doing stuff right now", 0, 11, 11},
		{"This is me, doing stuff right now", 12, 25, 23},
		{"Try it on my Rolento.", 3, 7, 6},
	}
	for _, c := range cases {
		got, _ := indexOfLastBreakChar(c.s, c.start, c.end)
		if got != c.want {
			t.Errorf("indexOfLastBreakChar(%v,%v,%v) == %v, want %v",
				c.s, c.start, c.end, got, c.want)
		}
	}
}
