package sqlick

import (
	"strings"
)

type SqlickColumn struct {
	GoName string
	GoType string
}

func (sc *SqlickColumn) String() string {
	return sc.GoName + ": " + sc.GoType
}

type SqlickTable struct {
	GoName  string
	Columns []SqlickColumn
}

func (sc *SqlickTable) String() string {
	sb := strings.Builder{}

	sb.WriteString(sc.GoName)
	for _, col := range sc.Columns {
		sb.WriteRune('\n')
		sb.WriteString("  ")
		sb.WriteString(col.String())
	}

	return sb.String()
}

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n")
}
