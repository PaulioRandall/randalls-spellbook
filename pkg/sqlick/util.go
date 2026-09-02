package sqlick

import (
	"fmt"
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

func genList[T any](
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
