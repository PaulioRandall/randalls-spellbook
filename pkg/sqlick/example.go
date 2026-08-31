package sqlick

import (
	"log"
)

type CheeseMaker struct {
	Id      int
	Name    string
	Country string
}

type Cheese struct {
	Id       int
	MakerId  int // Map to CheeseMaker.Id
	Name     string
	Strength int     // 1: Mild, 4: Extra Mature
	Rating   float64 // Out of 5
	Notes    string  // E.g. "Nutty after taste"
}

func sqlickExample() {
	path := "./testproject/sqlick.sqlite"
	db := NewSqliteDB(path)

	db.AddStructTable(CheeseMaker{})
	db.AddStructTable(Cheese{})

	e := db.Open()
	if e != nil {
		log.Fatal(e)
	}

	e = db.CreateTables()
	if e != nil {
		log.Fatal(e)
	}

	defer db.Close()

	// IDEA: After doing a SELECT, check if model implements
	//       a Init() function, and call it if it does.
	//       Init() can do things like populate private
	//       fields based on fetched data.
}
