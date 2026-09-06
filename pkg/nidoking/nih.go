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
		"  Needle: '%s',",
		"  Haystack: <not printed>,",
		"  LineNum(): %d,",
		"  InlineStart(): %d,",
		"  InlineEnd(): %d,",
		"  RuneLineStart(): %d,",
		"  RuneLineEnd(): %d,",
		"  RuneStart(): %d,",
		"  RuneEnd(): %d,",
		"  RuneInlineStart(): %d,",
		"  RuneInlineEnd(): %d,",
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
		nih.RuneLineStart(),
		nih.RuneLineEnd(),
		nih.RuneStart(),
		nih.RuneEnd(),
		nih.RuneInlineStart(),
		nih.RuneInlineEnd(),
	)
}

// LineNum returns the line number, i.e. LineIndex + 1.
func (nih NeedleInHaystack) LineNum() int {
	return nih.LineIndex + 1
}

// LineText returns the text of the whole line the needle
// was found on.
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

// RuneLineStart returns the rune index of the start of
// the line that the needle was found on.
func (nih NeedleInHaystack) RuneLineStart() int {
	return nih.runeIndex(nih.LineStart)
}

// RuneLineEnd returns the rune index of the end of
// the line that the needle was found on.
func (nih NeedleInHaystack) RuneLineEnd() int {
	return nih.runeIndex(nih.LineEnd)
}

// RuneStart returns the rune index of the needle's start.
func (nih NeedleInHaystack) RuneStart() int {
	return nih.runeIndex(nih.Start)
}

// RuneEnd returns the rune index of the needle's end.
func (nih NeedleInHaystack) RuneEnd() int {
	return nih.runeIndex(nih.End)
}

// RuneInlineStart returns the rune index of the needle's
// start relative to the start of the line it is on.
func (nih NeedleInHaystack) RuneInlineStart() int {
	return nih.runeIndexLine(nih.InlineStart())
}

// RuneInlineEnd returns the rune index of the needle's
// end relative to the start of the line it is on.
func (nih NeedleInHaystack) RuneInlineEnd() int {
	return nih.runeIndexLine(nih.InlineEnd())
}

// Replace replaces the instance of the needle with text
// and returns the new haystack.
func (nih NeedleInHaystack) Replace(
	text string,
) string {
	hay := nih.Haystack
	return hay[:nih.Start] + text + hay[nih.End:]
}

// ReplaceLine replaces the whole line the needle was found
// on with text and returns the updated haystack.
func (nih NeedleInHaystack) ReplaceLine(
	text string,
) string {
	hay := nih.Haystack
	return hay[:nih.LineStart] + text + hay[nih.LineEnd:]
}

// ReplaceList makes a copy of the line and replaces the
// needle for each string in texts. The original line is
// removed.
func (nih NeedleInHaystack) ReplaceList(
	texts []string,
) string {
	lines := make([]string, len(texts), len(texts))
	inStart := nih.InlineStart()
	inEnd := nih.InlineEnd()
	line := nih.LineText()

	for i, s := range texts {
		lines[i] = line[:inStart] + s + line[inEnd:]
	}

	newValue := strings.Join(lines, "\n")
	hay := nih.Haystack
	return hay[:nih.LineStart] + newValue + hay[nih.LineEnd:]
}

// RemoveLine removes the whole line the needle was found
// on and returns the updated haystack.
func (nih NeedleInHaystack) RemoveLine() string {
	hay := nih.Haystack

	if nih.LineStart > 0 {
		// Not first line.
		return hay[:nih.LineStart-1] + hay[nih.LineEnd:]
	}

	if len(hay) > nih.LineEnd {
		// But not last line.
		return hay[:nih.LineStart] + hay[nih.LineEnd+1:]
	}

	// Only one line in haystack.
	return ""
}

// FindNext finds the next instance of the needle in the
// haystack.
func (nih NeedleInHaystack) FindNext() NeedleInHaystack {
	return FindNeedle(nih.Haystack, nih.Needle, nih.End)
}

// ReplaceFindNext does the same as
// [NeedleInHaystack.Replace] but finds the next
// instance of the needle within the haystack.
func (nih NeedleInHaystack) ReplaceFindNext(
	text string,
) NeedleInHaystack {
	hay := nih.Replace(text)
	startingAt := nih.Start + len(text)
	return FindNeedle(hay, nih.Needle, startingAt)
}

// ReplaceLineFindNext does the same as
// [NeedleInHaystack.ReplaceFindNext] but for the
// line instead of the needle.
func (nih NeedleInHaystack) ReplaceLineFindNext(
	text string,
) NeedleInHaystack {
	hay := nih.ReplaceLine(text)
	startingAt := nih.Start + len(text)
	return FindNeedle(hay, nih.Needle, startingAt)
}

func (nih NeedleInHaystack) runeIndex(end int) int {
	return len([]rune(nih.Haystack[:end]))
}

func (nih NeedleInHaystack) runeIndexLine(end int) int {
	return len([]rune(
		nih.Haystack[nih.LineStart : nih.LineStart+end],
	))
}
