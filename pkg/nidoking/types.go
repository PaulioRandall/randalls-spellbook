package nidoking

import (
	"fmt"
	"strings"
)

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func trimLines(s string) string {
	lines := strings.Split(s, "\n")
	for i, v := range lines {
		lines[i] = strings.TrimSpace(v)
	}
	return strings.Join(lines, "\n")
}

// NeedleInHaystack holds byte position information about
// a substring (needle) within a bigger string (haystack).
// It is returned by [FindNeedle], and provides functions
// for accessing rune and line positions, and a few
// additional functions for common activities.
type NeedleInHaystack struct {
	// Line index.
	LineIndex int

	// Byte index of the first character of the line.
	LineStart int

	// Byte index one past the last character of the line
	// (exclusive, not including '\n').
	LineEnd int

	// Byte index of the first character of the search term.
	Start int

	// Byte index one past the last character of the search
	// term (exclusive).
	End int

	// The substring.
	Needle string

	// The string containing the substring.
	Haystack string
}

// String returns a developer understandable string
// representation of the object for debugging.
func (nih NeedleInHaystack) String() string {
	msg := joinLines(
		"NeedleInHaystack{",
		"  LineIndex: %d,",
		"  LineStart: %d,",
		"  LineEnd: %d,",
		"  Start: %d,",
		"  End: %d,",
		"  Needle: %s,",
		"  Haystack: <not printed>,",
		"  LineNum(): %d,",
		"  InlineStart(): %d,",
		"  InlineEnd(): %d,",
		"}",
	)

	return fmt.Sprintf(
		msg,
		nih.LineIndex,
		nih.LineStart,
		nih.LineEnd,
		nih.Start,
		nih.End,
		nih.Needle,
		nih.LineNum(),
		nih.InlineStart(),
		nih.InlineEnd(),
	)
}

// LineNum returns the line number, i.e. line index + 1.
func (nih NeedleInHaystack) LineNum() int {
	return nih.LineIndex + 1
}

// LineText returns the text of the whole line the needle
// is on.
func (nih NeedleInHaystack) LineText() string {
	return nih.Haystack[nih.LineStart:nih.LineEnd]
}

// InlineStart returns the starting byte index of the
// needle relative to the start of its line.
func (nih NeedleInHaystack) InlineStart() int {
	return nih.Start - nih.LineStart
}

// InlineEnd returns the ending byte index of the needle
// relative to the start of its line.
func (nih NeedleInHaystack) InlineEnd() int {
	return nih.End - nih.LineStart
}

// ReplaceNeedle replaces the instance of the needle within
// the haystack with text and returns the new haystack. The
// object itself will not be modified.
func (nih NeedleInHaystack) ReplaceNeedle(
	text string,
) string {
	hay := nih.Haystack
	return hay[:nih.Start] + text + hay[nih.End:]
}

// ReplaceLine replaces the whole line the needle was found
// on within the haystack with text and returns the updated
// haystack. The object itself will not be modified.
func (nih NeedleInHaystack) ReplaceLine(
	text string,
) string {
	hay := nih.Haystack
	return hay[:nih.LineStart] + text + hay[nih.LineEnd:]
}

// FindNext finds the next instance of the needle in the
// haystack.
func (nih NeedleInHaystack) FindNext() NeedleInHaystack {
	return FindNeedle(nih.Haystack, nih.Needle, nih.End)
}

// ReplaceNeedleFindNext does the same as
// [NeedleInHaystack.ReplaceNeedle] but finds the next
// instance of the needle wtihin the haystack.
func (nih NeedleInHaystack) ReplaceNeedleFindNext(
	text string,
) NeedleInHaystack {
	hay := nih.ReplaceNeedle(text)
	startingAt := nih.Start + len(text)
	return FindNeedle(hay, nih.Needle, startingAt)
}

// ReplaceLineFindNext does the same as
// [NeedleInHaystack.ReplaceNeedleFindNext] but for the
// line instead of the needle.
func (nih NeedleInHaystack) ReplaceLineFindNext(
	text string,
) NeedleInHaystack {
	hay := nih.ReplaceLine(text)
	startingAt := nih.Start + len(text)
	return FindNeedle(hay, nih.Needle, startingAt)
}
