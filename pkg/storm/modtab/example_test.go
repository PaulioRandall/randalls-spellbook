package modtab

import (
	"database/sql"
	"fmt"
)

func Example() {
	// YUDO: Error handling.

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	db, _ := sql.Open("sqlite", ":memory:")
	defer db.Close()

	// Create a model of the Player struct.
	model, tableExists, err := Map(db, Player{})

	// Create the Player table.
	err = model.Create(db)

	// Insert an entry into the Player table.
	// Can also be used to update an entry.
	err = model.Upsert(db, Player{
		Id:     123,
		Name:   "Bob",
		Rating: 6.9,
	})

	// Select all entries from the Player table.
	results, err := model.SelectAll[Player](db)

	// Select the entry from the Player table.
	result, found, err := model.SelectById[Player](db, 123)

	// Delete the entry from the Player table.
	err = model.DeleteById(db, 123)

	// Remove the Player table.
	err = model.Drop(db)

	_ = err
	_ = tableExists
	_ = found
	_ = results
	_ = result
}

func ExampleParse() {
	// YUDO: Error handling.

	type Player struct {
		ignored bool
		Id      int
		Name    string
		Rating  float64
	}

	model, _ := Parse(Player{})

	fmt.Printf("Struct: %s => Table: %s\n", model.GoName, model.SqlName)
	for _, prop := range model.Props {
		if prop.IsKey {
			fmt.Print("\tKey ")
		} else {
			fmt.Print("\t")
		}

		fmt.Printf(
			"Field: %s (%s) => Column: %s (%s)\n",
			prop.GoName,
			prop.GoType.Name(),
			prop.SqlName,
			prop.SqlType,
		)
	}
	// Output:
	// Struct: Player => Table: Player
	//	Key Field: Id (int) => Column: Id (INTEGER)
	//	Field: Name (string) => Column: Name (TEXT)
	//	Field: Rating (float64) => Column: Rating (REAL)
}

// The Map function excludes the Player.Rating field since
// it's not in the table, and correctly maps the tables
// primary key column to the Id field.
func ExampleMap() {
	// YUDO: Error handling.

	type Player struct {
		Name   string
		Rating float64
		Id     int
	}

	db, _ := sql.Open("sqlite", ":memory:")
	defer db.Close()

	_, _ = db.Exec(`
	CREATE TABLE Player (
		Id INTEGER NOT NULL DEFAULT 0,
		Name TEXT NOT NULL DEFAULT '',
	  PRIMARY KEY (Id)
	)
`)

	model, exists, _ := Map(db, Player{})
	fmt.Printf("Does table exist in database: %v\n", exists)

	fmt.Printf("Struct: %s => Table: %s\n", model.GoName, model.SqlName)
	for _, prop := range model.Props {
		if prop.IsKey {
			fmt.Print("\tKey ")
		} else {
			fmt.Print("\t")
		}

		fmt.Printf(
			"Field: %s (%s) => Column: %s (%s)\n",
			prop.GoName,
			prop.GoType.Name(),
			prop.SqlName,
			prop.SqlType,
		)
	}
	// Output:
	// Does table exist in database: true
	// Struct: Player => Table: Player
	//	Field: Name (string) => Column: Name (TEXT)
	//	Key Field: Id (int) => Column: Id (INTEGER)
}

func ExampleCachedMapper() {
	// YUDO: Error handling.

	type Player struct {
		Name   string
		Rating float64
		Id     int
	}

	db, _ := sql.Open("sqlite", ":memory:")
	defer db.Close()

	mapper := CachedMapper{}

	// Won't add the model to the cache becasue the table
	// doesn't exist.
	model, _, _ := mapper.Map(db, Player{})

	model.Create(db)

	// Will add the model to the cache because the table now
	// exists.
	model, _, _ = mapper.Map(db, Player{})

	model.Upsert(db, Player{
		Id:     123,
		Name:   "Bob",
		Rating: 6.9,
	})

	// Will pull from the cache rather than parsing again.
	model, _, _ = mapper.Map(db, Player{})

	// The table will be dropped but the cache entry still
	// remains.
	_ = model.Drop(db)

	// Removes the cache entry for the Player type.
	mapper.ClearType(Player{})

	// Can also remove all cache entries for the Player
	// table.
	mapper.ClearTable("Player")

	// Can also empty the entire cache.
	mapper.Clear()
}
