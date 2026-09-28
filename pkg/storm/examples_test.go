package storm

import (
// "fmt"
)

func ExampleStorm_Create() {
	// YUDO: Error handling.

	type Player struct {
		Id      int
		Name    string
		Rating  float64
		ignored *int
	}

	// Using in-memory database for the example but is
	// usually a local file path.
	db := New(":memory:")

	_ = db.Open()
	defer db.Close()

	_ = db.Create(Player{})
}

/*
func ExampleStorm_Table() {
	// YUDO: Error handling.

	type Player struct {
		Id      int
		Name    string
		Rating  float64
		ignored *int
	}

	// Using in-memory database for the example but is
	// usually a local file path.
	db := New(":memory:")

	_ = db.Open()
	defer db.Close()
	_ = db.Create(Player{})

	tableMeta, _ := db.Table(Player{})

	fmt.Println(tableMeta.TableName)
	for _, col := range tableMeta.Columns {
		if col.PrimaryKey {
			fmt.Printf("\t%s %s PRIMARY KEY\n", col.Name, col.Type)
		} else {
			fmt.Printf("\t%s %s\n", col.Name, col.Type)
		}
	}
	// Output:
	// Player
	//	Id INTEGER PRIMARY KEY
	//	Name TEXT
	//	Rating REAL
}
*/
