package storm

import (
	"log"
)

func Example() {
	var idPool int64 = 0
	generateId := func() int64 {
		idPool++
		return idPool
	}

	logFatal := func(e error) {
		if e != nil {
			log.Fatal(e)
		}
	}

	type Role struct {
		Id        int64
		Name      string
		Strength  int64
		Stamina   int64
		Intellect int64
		Health    int64
		Mana      int64
	}

	var guardian = Role{
		Id:        generateId(),
		Name:      "Guardian",
		Strength:  8,
		Stamina:   8,
		Intellect: 2,
		Health:    1000,
		Mana:      200,
	}

	var mage = Role{
		Id:        generateId(),
		Name:      "Mage",
		Strength:  3,
		Stamina:   6,
		Intellect: 10,
		Health:    500,
		Mana:      1000,
	}

	type Player struct {
		Id     int64
		Name   string
		RoleId int64
		role   *Role // This field is ignored.
	}

	var alice = Player{
		Id:     generateId(),
		Name:   "Alice",
		RoleId: mage.Id,
		role:   &mage,
	}

	var bob = Player{
		Id:     generateId(),
		Name:   "Bob",
		RoleId: guardian.Id,
		role:   &guardian,
	}

	var charlie = Player{
		Id:     generateId(),
		Name:   "Charlie",
		RoleId: mage.Id,
		role:   &mage,
	}

	path := "./game-files.sqlite"
	db := New(path)

	e := db.Open()
	logFatal(e)
	defer db.Close()

	// Create
	e = db.Create(Player{}, Role{})
	logFatal(e)

	// Insert
	e = db.Insert[any](guardian, mage, alice, bob, charlie)
	logFatal(e)

	// Update
	alice.Name = "Alicia"
	e = db.Update(alice)
	logFatal(e)

	// Select
	players, e := db.Select(Player{})
	logFatal(e)

	// Select (by ID)
	role, e := db.SelectById(Role{}, mage.Id)
	logFatal(e)

	_ = players
	_ = role
}
