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
	MakerId  int // Maps to CheeseMaker.Id
	Name     string
	Strength int     // 1: Mild, 4: Extra Mature
	Rating   float64 // Out of 5
	Notes    string  // E.g. "Nutty after taste"
}

var generateMakerId = newIntGenerator()
var generateCheeseId = newIntGenerator()

func newIntGenerator() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

func logFatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func sqlickExample() {
	path := "./testproject/sqlick.sqlite"
	db := NewSqliteDB(path)

	db.AddStructTable(CheeseMaker{})
	db.AddStructTable(Cheese{})

	e := db.Open()
	logFatal(e)
	defer db.Close()

	e = db.CreateTables()
	logFatal(e)

	charliesCheeses := CheeseMaker{
		Id:      generateMakerId(),
		Name:    "Charlie's Cheeses",
		Country: "England",
	}
	e = db.Insert(charliesCheeses)
	logFatal(e)

	charliesChedder := Cheese{
		Id:       generateCheeseId(),
		MakerId:  charliesCheeses.Id,
		Name:     "Charlie's Chedder",
		Strength: 1,
		Rating:   3.5,
		Notes:    "Charlie's standard mild Chedder cheese",
	}
	e = db.Insert(charliesChedder)
	logFatal(e)

	charliesMatureChedder := Cheese{
		Id:       generateCheeseId(),
		MakerId:  charliesCheeses.Id,
		Name:     "Charlie's Mature Chedder",
		Strength: 4,
		Rating:   4.8,
	}
	e = db.Insert(charliesMatureChedder)
	logFatal(e)

	vanillaCharm := Cheese{
		Id:       generateCheeseId(),
		MakerId:  charliesCheeses.Id,
		Name:     "Vanilla Charm",
		Strength: 2,
		Rating:   4.3,
		Notes:    "Sweeter than typical Chedder with strong vanilla tones",
	}
	e = db.Insert(vanillaCharm)
	logFatal(e)

	francsSheepFromage := CheeseMaker{
		Id:      generateMakerId(),
		Name:    "Franc's Sheep Fromage",
		Country: "France",
	}
	e = db.Insert(francsSheepFromage)
	logFatal(e)

	blueSheep := Cheese{
		Id:       generateCheeseId(),
		MakerId:  francsSheepFromage.Id,
		Name:     "Blue Sheep",
		Strength: 2,
		Rating:   4.4,
		Notes:    "Tangy Roquefort cheese",
	}
	e = db.Insert(blueSheep)
	logFatal(e)

	sardaSheep := Cheese{
		Id:       generateCheeseId(),
		MakerId:  francsSheepFromage.Id,
		Name:     "Sarda Sheep",
		Strength: 1,
		Rating:   3.3,
		Notes:    "Hard sheep cheese from Sarda, Italy",
	}
	e = db.Insert(sardaSheep)
	logFatal(e)

	// IDEA: After doing a SELECT, check if model implements
	//       a Init() function, and call it if it does.
	//       Init() can do things like populate private
	//       fields based on fetched data.
}
