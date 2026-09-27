package storm

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/PaulioRandall/randalls-spellbook/pkg/curse"
)

var (
	ErrMkDirPath = curse.Err(
		"Unable to check or create path to SQLite database",
	)
)

func typeName(model any) string {
	return reflect.TypeOf(model).Name()
}

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n")
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

	return ErrMkDirPath.Wrap(e)
}
