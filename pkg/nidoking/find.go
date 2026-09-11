package nidoking

import (
	"regexp"
	"strings"
)

const (
	// ModeString is set as [NeedleInHaystack].Mode for
	// the results of [Find] and results of subsequent calls
	// to [Replacement.FindNext] and
	// [NeedleInHaystack.FindNext].
	ModeString = "string"

	// ModeString is set as [NeedleInHaystack].Mode for
	// the results of [Match] and results of subsequent calls
	// to [Replacement.FindNext] and
	// [NeedleInHaystack.FindNext].
	ModeRegexp = "regexp"
)

// Find locates needle within haystack and returns a
// NeedleInHaystack. If the needle contains any linefeeds
// then panic ensues. If no match is found then an empty
// NeedleInHaystack is returned; use
// [NeedleInHaystack.IsMatch] to avoid referencing the
// return type.
func Find(
	haystack, needle string,
	startingAt int,
) NeedleInHaystack {
	panicIfMultilineNeedle(needle)

	start := strings.Index(haystack[startingAt:], needle)
	if start == -1 {
		return NeedleInHaystack{}
	}

	// Correct for substring of haystack used.
	start += startingAt
	end := start + len(needle)

	lineStart, lineEnd := findLineIndexes(haystack, start)

	return NeedleInHaystack{
		Start:     start,
		End:       end,
		LineStart: lineStart,
		LineEnd:   lineEnd,
		Mode:      ModeString,
		Needle:    needle,
		Pattern:   needle,
		Haystack:  haystack,
	}
}

// Find locates needle within haystack and returns a
// NeedleInHaystack. If the needle contains any linefeeds
// then panic ensues. If no match is found then an empty
// NeedleInHaystack is returned; use
// [NeedleInHaystack.IsMatch] to avoid referencing the
// return type.
func Match(
	haystack, pattern string,
	startingAt int,
) NeedleInHaystack {
	panicIfMultilinePattern(pattern)

	re := regexp.MustCompile(pattern)
	pos := re.FindStringIndex(haystack[startingAt:])
	if pos == nil {
		return NeedleInHaystack{}
	}

	start := pos[0] + startingAt
	end := pos[1] + startingAt

	lineStart, lineEnd := findLineIndexes(haystack, start)

	return NeedleInHaystack{
		Start:     start,
		End:       end,
		LineStart: lineStart,
		LineEnd:   lineEnd,
		Mode:      ModeRegexp,
		Needle:    haystack[start:end],
		Pattern:   pattern,
		Haystack:  haystack,
	}
}

func findLineIndexes(
	haystack string,
	start int,
) (int, int) {
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

	return lineStart, lineEnd
}

func panicIfMultilineNeedle(needle string) {
	if strings.ContainsRune(needle, '\n') {
		panic("Needle must not be a multi-line string: \\n")
	}

	if strings.ContainsRune(needle, '\r') {
		panic("Needle must not be a multi-line string: \\r")
	}
}

func panicIfMultilinePattern(pattern string) {
	if strings.Contains(pattern, "\\n") {
		panic("Pattern must not be a multi-line string: \\n")
	}

	if strings.Contains(pattern, "\\r") {
		panic("Pattern must not be a multi-line string: \\r")
	}

	// matches (?m), (?s), (?im), etc.
	// I assumed flags are case-sensitive.
	re := regexp.MustCompile(`\n|\r|\(\?[iU]*[ms].*\)`)
	if s := re.FindString(pattern); s != "" {
		panic("Pattern must not be a multi-line string: " + s)
	}
}
