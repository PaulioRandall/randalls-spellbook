// Package storm is a minimalistic ORM based database
// API for storing and accessing simple structures/tables.
//
// Instead of writing or building query statements, users
// design structs with the storage structure they want.
// Objects with these struct types are passed to various
// functions to perform standard database actions such as
// creating tables and inserting data. It's an approach
// used by existing database packages, such as
// https://gorm.io/. This package provides a minimalistic
// solution that trades away the feature richness of
// existing tools for ease of use when developing simple
// desktop tools and while rapid prototyping. It also aims
// for high plunderability, i.e. the ability for people to
// copy and modify the codebase for their own purposes.
//
// Throughout this package the term 'model' refers to an
// object passed purely for its type (create table, etc),
// while 'object' is used when the data is the subject
// (insert, update, etc).
//
//	type Person struct {
//		Id int64 // First column is always the primary key.
//		Name string
//		Height float64 // In meters.
//		isTall bool // Unexported fields are ignored.
//	}
//
//	db, err := NewQlite("./db.sqlite")
//	// YUDO: Handle error.
//
//	err = db.Create(Person{})
//	// YUDO: Handle error.
//
//	err = db.Open()
//	// YUDO: Handle error.
//	defer db.Close()
//
//	alice := Person{
//		Id: 1,
//		Name: "Alice",
//		Height: 1.59,
//		isTall: false,
//	}
//	err = db.Insert(alice)
//	// YUDO: Handle error.
//
//	bob := Person{
//		Id: 2,
//		Name: "Bob",
//		Height: 1.92,
//		isTall: true,
//	}
//	err = db.Insert(bob)
//	// YUDO: Handle error.
//
//	alice.Height = 1.6
//	err = db.Update(alice)
//	// YUDO: Handle error.
//
//	people, err := db.SelectAll(Person{})
//	// YUDO: Handle error.
//
//	bob, err = db.SelectById(Person{}, 2)
//	// YUDO: Handle error.
//
//	err = db.Delete(alice)
//	// YUDO: Handle error.
//
//	err = db.Drop(Person{})
//	// YUDO: Handle error.
package storm
