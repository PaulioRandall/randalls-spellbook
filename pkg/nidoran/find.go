package nidoran

import (
	"strings"
)

// Find locates needle within haystack and returns a
// NeedleInHaystack. If no match is found then an empty
// NeedleInHaystack is returned. If the needle contains
// any linefeeds then panic ensues.
func Find(
	haystack, needle string,
	startingAt int,
) NeedleInHaystack {
	if strings.Contains(needle, "\n") {
		panic("Needle must not be a multi-line string")
	}

	start := strings.Index(haystack[startingAt:], needle)
	if start == -1 {
		return NeedleInHaystack{}
	}

	// Correct for substring of haystack used.
	start += startingAt
	end := start + len(needle)

	// Find start of line: scan backwards for the previous
	// '\n'. +1 because we want to exclude the linefeed.
	lineStart := strings.LastIndexByte(
		haystack[:start],
		'\n',
	) + 1

	if lineStart == -1 {
		lineStart = 0
	}

	// Find end of line: scan forwards for the next '\n'.
	lineEnd := strings.IndexByte(haystack[start:], '\n')

	if lineEnd == -1 {
		lineEnd = len(haystack)
	} else {
		// Correct for substring of haystack used.
		lineEnd += start
	}

	lineIdx := strings.Count(haystack[:start], "\n")

	return NeedleInHaystack{
		LineIndex: lineIdx,
		LineStart: lineStart,
		LineEnd:   lineEnd,
		Start:     start,
		End:       end,
		Needle:    needle,
		Haystack:  haystack,
	}
}
