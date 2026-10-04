package storm

import (
	"fmt"
)

func ExampleStorm_Table() {
	// YUDO: Error handling.

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	st, err := Open(":memory:")
	defer st.Close()

	err = st.Create(Player{})
	table, err := st.Table(Player{})
	_ = err

	fmt.Println(table.Name)
	for _, col := range table.Columns {
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

func ExampleStorm_Create() {
	// YUDO: Error handling.

	type Player struct {
		Id      int
		Name    string
		Rating  float64
		ignored *int // Will be ignored.
	}

	// Using in-memory database for the example but is
	// usually a local file path.
	st, err := Open(":memory:")
	defer st.Close()

	err = st.Create(Player{})
	_ = err
}

func ExampleStorm_Drop() {
	// YUDO: Error handling.

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	st, err := Open(":memory:")
	defer st.Close()

	err = st.Create(Player{})
	err = st.Drop(Player{})
	_ = err
}

func ExampleStorm_Put() {
	// YUDO: Error handling.

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	st, err := Open(":memory:")
	defer st.Close()

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

	err = st.Put(players...)
	_ = err
}

func ExampleStorm_List() {
	// YUDO: Error handling.

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	st, err := Open(":memory:")
	defer st.Close()

	// Table auto created on insert.
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

	players, err := st.List(Player{}, "Rating > ?", 2)
	_ = err

	for _, p := range players {
		fmt.Printf("%d: %s\n", p.Id, p.Name)
	}
	// Output:
	// 2: Bob
	// 3: Charlie
}

func ExampleStorm_Get() {
	// YUDO: Error handling.

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	st, err := Open(":memory:")
	defer st.Close()

	// Table auto created on insert.
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

	bob, err := st.Get(Player{}, "Id = ?", 2)
	_ = err

	fmt.Printf("%d: %s\n", bob.Id, bob.Name)
	// Output:
	// 2: Bob
}

func ExampleStorm_Delete() {
	// YUDO: Error handling.

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	st, err := Open(":memory:")
	defer st.Close()

	// Table auto created on insert.
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

	err = st.Delete(Player{}, "Id = ?", 2)
	_ = err

	players, err := st.List(Player{}, "")
	_ = err

	for _, p := range players {
		fmt.Printf("%d: %s\n", p.Id, p.Name)
	}
	// Output:
	// 1: Alice
	// 3: Charlie
}
