package schema

import (
	"database/sql"
	"errors"
)

type Table struct {
	Name    string
	Columns []Column
}

type Column struct {
	Name       string
	Type       string
	PrimaryKey bool
}

func Query(db *sql.DB, tableName string) (Table, error) {
	var t Table

	tableSchema, e := QueryTable(db, tableName)

	if errors.Is(e, ErrTableNotFound) {
		// Not found
		return t, nil
	}

	if e != nil {
		return t, e
	}

	colSchema, e := QueryColumns(db, tableName)
	if e != nil {
		return t, e
	}

	t.Name = tableSchema.TableName
	t.Columns = make([]Column, len(colSchema), len(colSchema))

	for i, col := range colSchema {
		t.Columns[i] = Column{
			Name:       col.Name,
			Type:       col.Type,
			PrimaryKey: col.PrimaryKey > 0,
		}
	}

	return t, nil
}
