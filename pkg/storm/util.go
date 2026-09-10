package storm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type fmtBuilder struct {
	strings.Builder
}

func (fb *fmtBuilder) WriteFmt(msg string, args ...any) {
	s := fmt.Sprintf(msg, args...)
	fb.WriteString(s)
}

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func generateList[T any](
	list []T,
	itemToStr func(int, T) (string, error),
) (string, error) {
	fb := fmtBuilder{}

	for i, item := range list {
		if i != 0 {
			fb.WriteString(",\n")
		}

		s, e := itemToStr(i, item)
		if e != nil {
			return "", e
		}

		fb.WriteString(s)
	}

	return fb.String(), nil
}

func makeParentDirs(path string) error {
	if path == ":memory" {
		// SQlite in-memory database. There is no path!
		return nil
	}

	parent := filepath.Dir(path)
	e := os.MkdirAll(parent, os.ModePerm)
	if e == nil {
		return nil
	}

	return fmt.Errorf(
		"Unable to check or create directory path to SQLite database: %w",
		e,
	)
}

func errOrNil(e error, msg string, args ...any) error {
	if e != nil {
		return errMaybeWrap(e, msg, args...)
	}
	return nil
}

func errMaybeWrap(e error, msg string, args ...any) error {
	if e == nil {
		return fmt.Errorf(msg, args...)
	}
	return fmt.Errorf(msg+": %w", append(args, e)...)
}
