package nidoran

import (
	"fmt"
	"strings"
)

type stringer interface {
	String() string
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

func stringifyValue[T any](value T) string {
	if st, ok := any(value).(stringer); ok {
		return st.String()
	}

	return fmt.Sprintf("%v", value)
}

func stringifyValues[T any](values []T) []string {
	strs := make([]string, len(values), len(values))

	for i, v := range values {
		strs[i] = stringifyValue(v)
	}

	return strs
}

func fmtPrintString(s string) string {
	const maxStringLength = 30
	s = strings.Replace(s, "\n", "\\n", -1)

	r := []rune(s)
	if len(r) > maxStringLength {
		s = string(r[:maxStringLength]) + "..."
	}

	return string(s)
}
