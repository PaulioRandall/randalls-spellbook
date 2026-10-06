package stormy

import (
	"fmt"

	"github.com/PaulioRandall/randalls-spellbook/pkg/wizzard"
)

func ExampleOpenMode() {
	// YUDO: Error handling.

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	st := New(":memory:")
	defer st.Close()

	fmt.Printf("Before any operation, IsOpen: %v\n", st.IsOpen())

	// Stormy defaults to OpenModePersist which opens the
	// database the moment you try to perform an operation
	// then remains open until manually closed.
	err := st.Create(Player{})
	fmt.Printf("OpenModePersist, IsOpen: %v\n", st.IsOpen())
	st.Close()

	// OpenModeRequest will open a database to perform an
	// operation then close. However, if the database is
	// already open it is left open.
	st.SetOpenMode(OpenModeRequest)
	err = st.Create(Player{})
	fmt.Printf("OpenModeRequest, IsOpen: %v\n", st.IsOpen())

	// OpenModeManual disables auto opening and closing
	// entirely. An error is returned if the database isn't
	// open.
	st.SetOpenMode(OpenModeManual)
	err = st.Create(Player{})
	fmt.Printf("OpenModeManual, Error: %v\n", err)
}

func ExampleStormy_Table() {
	// YUDO: Error handling.

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	st, err := Open(":memory:")
	defer st.Close()
	err = st.Create(Player{})

	// Point of interest.
	model, err := st.Table(Player{})

	_ = err

	fmt.Println(model.SqlTableString())
	// Output:
	// Sql Table: Player
	//	Id INTEGER
	//	Name TEXT
	//	Rating REAL
}

func ExampleStormy_Create() {
	// YUDO: Error handling.

	type Player struct {
		Id      int
		Name    string
		Rating  float64
		ignored *int // Will be ignored.
	}

	st, err := Open(":memory:")
	defer st.Close()

	// Point of interest.
	err = st.Create(Player{})

	_ = err
}

func ExampleStormy_Drop() {
	// YUDO: Error handling.

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	st, err := Open(":memory:")
	defer st.Close()
	err = st.Create(Player{})

	// Point of interest.
	err = st.Drop(Player{})

	_ = err
}

func ExampleStormy_Put() {
	// YUDO: Error handling.

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	st, err := Open(":memory:")
	defer st.Close()

	// Point of interest.
	// Table auto created on insert.
	err = st.Put(Player{
		Id:     4,
		Name:   "Dave",
		Rating: 4.4,
	})

	players := []Player{
		Player{
			Id:     1,
			Name:   "Alice",
			Rating: 1.1,
		},
		Player{
			Id:     2,
			Name:   "Bob",
			Rating: 2.2,
		},
		Player{
			Id:     3,
			Name:   "Charlie",
			Rating: 3.3,
		},
	}

	// Point of interest.
	err = st.Put(players...)

	_ = err
}

func ExampleStormy_List() {
	// YUDO: Error handling.

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	st, err := Open(":memory:")
	defer st.Close()

	err = st.Put(
		Player{
			Id:     1,
			Name:   "Alice",
			Rating: 1.1,
		},
		Player{
			Id:     2,
			Name:   "Bob",
			Rating: 2.2,
		},
		Player{
			Id:     3,
			Name:   "Charlie",
			Rating: 3.3,
		},
	)

	// Point of interest.
	players, err := st.List(Player{}, "Rating > ?", 2)

	_ = err

	for _, p := range players {
		fmt.Printf("%d: %s\n", p.Id, p.Name)
	}
	// Output:
	// 2: Bob
	// 3: Charlie
}

func ExampleStormy_Get() {
	// YUDO: Error handling.

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	st, err := Open(":memory:")
	defer st.Close()

	err = st.Put(
		Player{
			Id:     1,
			Name:   "Alice",
			Rating: 1.1,
		},
		Player{
			Id:     2,
			Name:   "Bob",
			Rating: 2.2,
		},
		Player{
			Id:     3,
			Name:   "Charlie",
			Rating: 3.3,
		},
	)

	// Point of interest.
	bob, found, err := st.Get(Player{}, "Id = ?", 2)

	_ = err

	fmt.Printf("Found: %v\n", found)
	fmt.Printf("%d: %s\n", bob.Id, bob.Name)
	// Output:
	// Found: true
	// 2: Bob
}

func ExampleStormy_Delete() {
	// YUDO: Error handling.

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	st, err := Open(":memory:")
	defer st.Close()

	err = st.Put(
		Player{
			Id:     1,
			Name:   "Alice",
			Rating: 1.1,
		},
		Player{
			Id:     2,
			Name:   "Bob",
			Rating: 2.2,
		},
		Player{
			Id:     3,
			Name:   "Charlie",
			Rating: 3.3,
		},
	)

	// Point of interest.
	err = st.Delete(Player{}, "Id = ?", 2)

	players, err := st.List(Player{}, "")
	_ = err

	for _, p := range players {
		fmt.Printf("%d: %s\n", p.Id, p.Name)
	}
	// Output:
	// 1: Alice
	// 3: Charlie
}

func ExampleStormy_Custom() {
	// YUDO: Error handling.

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	st, err := Open(":memory:")
	defer st.Close()

	err = st.PutAs(
		"User",
		Player{
			Id:     1,
			Name:   "Alice",
			Rating: 1.1,
		},
		Player{
			Id:     2,
			Name:   "Bob",
			Rating: 2.2,
		},
		Player{
			Id:     3,
			Name:   "Charlie",
			Rating: 3.3,
		},
	)

	// Point of interest.
	listLi := func(ctx OperationContext) ([]Player, error) {
		model, exist, err := ctx.MapAs("User", Player{})

		if !exist {
			return nil, nil
		}

		db := ctx.Database()
		rows, err := db.Query(`
		SELECT
			Id,
			Name,
			Rating
		FROM
			User
		WHERE Name LIKE '%li%'
	`)

		results, err := wizzard.ScanRows[Player](model, rows)
		_ = err
		return results, nil
	}

	// Point of interest.
	results, err := st.Custom(listLi)

	for _, p := range results {
		fmt.Printf("%d: %s\n", p.Id, p.Name)
	}

	_ = err
	// Output:
	// 1: Alice
	// 3: Charlie
}
