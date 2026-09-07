package nidoran

import (
	"fmt"
	"strings"
)

// Replacement is the result of a replacement function
// being called on [NeedleInHaystack].
type Replacement struct {
	// Nih is the original [NeedleInHaystack].
	Nih NeedleInHaystack

	// Start is the byte index of the start of the replaced
	// content within Haystack.
	Start int

	// End is the byte index of the end of the replaced
	// content within Haystack.
	End int

	// The result of the replacement.
	Haystack string
}

// FindNext finds the next instance of the needle in the
// haystack.
func (rep Replacement) FindNext() NeedleInHaystack {
	return Find(rep.Haystack, rep.Nih.Needle, rep.End)
}

// NeedleInHaystack holds byte position information about
// a substring (needle) within a bigger string (haystack).
// It is returned by [Find], and provides functions
// for accessing rune and line positions, and a few
// for replacing and removing the needle or its line.
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

// ReplaceInline replaces the instance of the needle with
// text and returns a [Replacement] object.
func (nih NeedleInHaystack) ReplaceInline(
	text string,
) Replacement {
	nih.panicIfEmpty()
	return nih.replacement(nih.Start, nih.End, text)
}

// ReplaceLine replaces the line the needle was found on
// with text and returns a [Replacement] object.
func (nih NeedleInHaystack) ReplaceLine(
	text string,
) Replacement {
	nih.panicIfEmpty()
	return nih.replacement(nih.LineStart, nih.LineEnd, text)
}

// ReplaceJoin copies the line and replaces the needle
// for each string in texts then adds the delim to the end
// of the line, except the last. The original line is
// removed. Returns a [Replacement] object.
func (nih NeedleInHaystack) ReplaceJoin(
	texts []string,
	delim string,
) Replacement {
	nih.panicIfEmpty()

	start := nih.InlineStart()
	end := nih.InlineEnd()
	line := nih.LineText()

	length := len(texts)
	lines := make([]string, length, length)

	for i, s := range texts {
		lines[i] = line[:start] + s + line[end:]
	}

	text := strings.Join(lines, delim+"\n")
	return nih.replacement(nih.LineStart, nih.LineEnd, text)
}

// ReplaceRepeat copies the line n times, replaces the
// needle with text, and adds the delim to the end of each
// line, except the last. The original line is removed.
// Returns a [Replacement] object.
func (nih NeedleInHaystack) ReplaceRepeat(
	text string,
	n int,
	delim string,
) Replacement {
	nih.panicIfEmpty()

	start := nih.InlineStart()
	end := nih.InlineEnd()
	line := nih.LineText()
	line = line[:start] + text + line[end:]

	lines := make([]string, n, n)

	for i := 0; i < n; i++ {
		lines[i] = line
	}

	s := strings.Join(lines, delim+"\n")
	return nih.replacement(nih.LineStart, nih.LineEnd, s)
}

// RemoveLine removes the whole line the needle was found
// on and returns the updated haystack. Lines cannot be
// removed using [NeedleInHaystack.ReplaceLine] function
// as it doesn't remove linefeeds.
func (nih NeedleInHaystack) RemoveLine() Replacement {
	nih.panicIfEmpty()

	rep := Replacement{
		Nih:   nih,
		Start: nih.LineStart,
		End:   nih.LineStart,
	}

	hay := nih.Haystack

	if nih.LineStart > 0 {
		// When not first line.
		rep.Haystack = hay[:nih.LineStart-1] + hay[nih.LineEnd:]
		return rep
	}

	if len(hay) > nih.LineEnd {
		// When not last line.
		rep.Haystack = hay[:nih.LineStart] + hay[nih.LineEnd+1:]
		return rep
	}

	// Single line haystack becomes an empty string.
	return rep
}

// FindNext finds the next instance of the needle in the
// haystack.
func (nih NeedleInHaystack) FindNext() NeedleInHaystack {
	nih.panicIfEmpty()
	return Find(nih.Haystack, nih.Needle, nih.End)
}

func (nih NeedleInHaystack) replacement(
	start, end int,
	text string,
) Replacement {
	return Replacement{
		Nih:   nih,
		Start: start,
		End:   start + len(text),
		Haystack: nih.Haystack[:start] +
			text +
			nih.Haystack[end:],
	}
}

func (nih NeedleInHaystack) runeIndex(end int) int {
	return len([]rune(nih.Haystack[:end]))
}

func (nih NeedleInHaystack) runeIndexLine(end int) int {
	return len([]rune(
		nih.Haystack[nih.LineStart : nih.LineStart+end],
	))
}

func (nih NeedleInHaystack) panicIfEmpty() {
	if nih == (NeedleInHaystack{}) {
		panic("Can't operate on empty NeedleInHaystack")
	}
}

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
