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
	Name    string
	Columns []DbColumn
}

type DbColumn struct {
	Name     string
	DataType string
}

func (table DbTable) String() string {
	sb := strings.Builder{}

	sb.WriteString(table.Name)

	for _, col := range table.Columns {
		sb.WriteString("\n  ")
		sb.WriteString(col.String())
	}

	return sb.String()
}

func (col DbColumn) String() string {
	sb := strings.Builder{}

	sb.WriteString(col.Name)
	sb.WriteString(": ")
	sb.WriteString(col.DataType)

	return sb.String()
}

func (ds *SqliteDatastore) Tables() []DbTable {
	return ds.tables
}

func (ds *SqliteDatastore) AddEntity(object any) error {
	table, e := parseDbTable(object)
	if e != nil {
		return e
	}

	ds.tables = append(ds.tables, table)
	return nil
}
