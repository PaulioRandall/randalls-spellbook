package data

import (
	"database/sql"
	"fmt"
)

var currentDataVersion int = 1

// Project holds general information about the project.
//
// This is the one table that never hopefully never
// experience breaking changes, although, additions are
// expected.
type Project struct {
	// DataVersion is the version of the data structure.
	DataVersion int
}

// GetDataVersion returns the version of the project data
// structure.
func (p *Project) GetDataVersion() int {
	return p.DataVersion
}

func createProjectTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS project (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			data_version INTEGER NOT NULL
		)
	`

	_, e := db.Exec(query)

	if e != nil {
		return fmt.Errorf("Failed to create project table: %w", e)
	}

	proj, e := getProject(db)
	if e != nil {
		return e
	}

	if proj == (Project{}) {
		return insertProject(db)
	}

	return nil
}

func insertProject(db *sql.DB) error {
	query := `
		INSERT INTO project (
			data_version
		) VALUES (?)
	`

	_, e := db.Exec(
		query,
		currentDataVersion,
	)

	if e != nil {
		return fmt.Errorf("Failed to insert project: %w", e)
	}

	return nil
}

func updateProject(db *sql.DB, project Project) error {
	query := `
		UPDATE
			project
		SET
			data_version = ?
	`

	_, e := db.Exec(query, project.DataVersion)

	if e != nil {
		return fmt.Errorf("Failed to update project table: %w", e)
	}

	return nil
}

func getProject(db *sql.DB) (Project, error) {
	query := `
		SELECT
			data_version
		FROM
			project
	`

	empty := Project{}

	rows, e := db.Query(query)
	if e != nil {
		return empty, fmt.Errorf("Failed to query project table: %w", e)
	}

	defer rows.Close()

	projects, e := parseProjectTableRows(rows)
	if e != nil {
		return empty, e
	}

	if len(projects) == 0 {
		return empty, nil
	}

	return projects[0], nil
}

func parseProjectTableRows(rows *sql.Rows) ([]Project, error) {
	var result []Project

	for rows.Next() {
		var p Project

		e := rows.Scan(
			&p.DataVersion,
		)

		if e != nil {
			return nil, fmt.Errorf("Failed to scan project table row: %w", e)
		}

		result = append(result, p)
	}

	return result, rows.Err()
}
