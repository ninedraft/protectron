package registry

import (
	"testing"
)

type User struct {
	ID    int    // primary key
	Group string `storm:"index"`  // this field will be indexed
	Email string `storm:"unique"` // this field will be indexed with a unique constraint
	Name  string // this field will not be indexed
	Age   int    `storm:"index"`
}

func TestDB(test *testing.T) {
	test.Log("creating db")
	var registry, err = New("testdata/test.db")
	if err != nil {
		test.Fatal("unable to open db file", err)
	}

	test.Log("creating user")
	var tx, errBegin = registry.db.Begin(true)
	if errBegin != nil {
		test.Fatal(errBegin)
	}
	defer tx.Rollback()
	if err := registry.db.Update(&User{
		Name: "merlin",
	}); err != nil {
		test.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		test.Fatal(err)
	}
}

func assertNilErr(test *testing.T, err error, message string) {
	if err != nil {
		test.Fatal(message, err)
	}
}
