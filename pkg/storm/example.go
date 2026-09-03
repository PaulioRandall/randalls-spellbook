package storm

import (
	"log"
)

type CheeseMaker struct {
	Id      int64
	Name    string
	Country string
}

type Cheese struct {
	Id       int64
	MakerId  int64 // Maps to CheeseMaker.Id
	Name     string
	Strength int64   // 1: Mild, 4: Extra Mature
	Rating   float64 // Out of 5
	Notes    string  // E.g. "Nutty after taste"
}

var generateMakerId = newIntGenerator()
var generateCheeseId = newIntGenerator()

var cheeseMakers = []CheeseMaker{
	CheeseMaker{
		Id:      generateMakerId(),
		Name:    "Charlie's Cheeses",
		Country: "England",
	},
	CheeseMaker{
		Id:      generateMakerId(),
		Name:    "Franc's Sheep Fromage",
		Country: "France",
	},
}

var cheeses = []Cheese{
	Cheese{
		Id:       generateCheeseId(),
		MakerId:  cheeseMakers[0].Id,
		Name:     "Charlie's Cheddar",
		Strength: 1,
		Rating:   3.5,
		Notes:    "Charlie's standard mild Cheddar cheese",
	},
	Cheese{
		Id:       generateCheeseId(),
		MakerId:  cheeseMakers[0].Id,
		Name:     "Charlie's Mature Cheddar",
		Strength: 4,
		Rating:   4.8,
	},
	Cheese{
		Id:       generateCheeseId(),
		MakerId:  cheeseMakers[0].Id,
		Name:     "Vanilla Charm",
		Strength: 2,
		Rating:   4.3,
		Notes:    "Sweeter than typical Cheddar with strong vanilla tones",
	},
	Cheese{
		Id:       generateCheeseId(),
		MakerId:  cheeseMakers[1].Id,
		Name:     "Blue Sheep",
		Strength: 2,
		Rating:   4.4,
		Notes:    "Tangy Roquefort cheese",
	},
	Cheese{
		Id:       generateCheeseId(),
		MakerId:  cheeseMakers[1].Id,
		Name:     "Sarda Sheep",
		Strength: 1,
		Rating:   3.3,
		Notes:    "Hard sheep cheese from Sarda, Italy",
	},
}

func newIntGenerator() func() int64 {
	var i int64 = 0
	return func() int64 {
		i++
		return i
	}
}

func logFatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func sqlick_example_1() {
	path := "./sqlick.sqlite"
	db := New(path)

	e := db.Open()
	logFatal(e)
	defer db.Close()

	e = db.Create(CheeseMaker{})
	logFatal(e)

	e = db.Create(Cheese{})
	logFatal(e)

	// Insert cheese makers and cheeses.
	for _, maker := range cheeseMakers {
		e = db.Insert(maker)
		logFatal(e)
	}
	for _, ch := range cheeses {
		e = db.Insert(ch)
		logFatal(e)
	}

	// Increase the rating of Charlie's Cheddar cheese then
	// update it in the database.
	ch := cheeses[0]
	ch.Rating += 0.3
	e = db.Update(ch)
	logFatal(e)

	var result any
	var cheeseList []Cheese
	var vanillaCharm Cheese

	result, e = db.SelectAll(Cheese{})
	logFatal(e)
	cheeseList, _ = result.([]Cheese)
	_ = cheeseList

	result, e = db.SelectById(Cheese{}, 3)
	logFatal(e)
	vanillaCharm, _ = result.(Cheese)
	_ = vanillaCharm

	/*
		// Select a specific cheese from the database. The second
		// argument must be the ID of the field.
		cheeseId := 1
		resultCheese, e = db.Select(Cheese{}, cheeseId)
		logFatal(e)
		cheeses, _ = resultCheese.(Cheese)

		// Select all cheese from the database.
		var cheeses []Cheese
		_, e := db.Select(&cheeses, nil)
		logFatal(e)

		// Select all cheeses from the database. This approach
		// requires explicit type casting.
		result, e := db.Select(Cheese{}, nil)
		logFatal(e)
		cheeses, _ = result.([]Cheese)
	*/

	// IDEA: After doing a SELECT, check if model implements
	//       a Init() function, and call it if it does.
	//       Init() can do things like populate private
	//       fields based on fetched data.

	// IDEA: Pass pointer to array to Select.
	// Select all cheeses from the database appending them
	// to the passed slice (note that a pointer to the slice
	// must be passed). The pointer to the array will be
	// returned but you'll usually want to ignore it.
	// var cheeses []Cheese
	// _, e := db.Select(&cheeses, nil)
	// logFatal(e)
}
