package datastore

import (
	"database/sql"
	"strings"

	_ "github.com/glebarez/go-sqlite"
)

type Datastore interface {
	Open() error
	Close() error
	AddEntity(object any) error
}

type SqliteDatastore struct {
	db     *sql.DB
	tables []DbTable
}

type DbTable struct {
	name    string
	columns []DbColumn
}

type DbColumn struct {
	name    string
	sqlType string
}

func (table DbTable) String() string {
	sb := strings.Builder{}

	sb.WriteString(table.name)

	for _, col := range table.columns {
		sb.WriteString("\n  ")
		sb.WriteString(col.String())
	}

	return sb.String()
}

func (col DbColumn) String() string {
	sb := strings.Builder{}

	sb.WriteString(col.name)
	sb.WriteString(": ")
	sb.WriteString(col.sqlType)

	return sb.String()
}

func (ds *SqliteDatastore) Tables() []DbTable {
	return ds.tables
}
