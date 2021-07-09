package test

import (
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/hlxwell/gorm-by-example/db"
)

var Fixtures *testfixtures.Loader

func LoadFixtures() {
	configFixtures()
	if err := Fixtures.Load(); err != nil {
		panic(err)
	}
}

func configFixtures() {
	conn, err := db.Conn.DB()
	if err != nil {
		panic(err)
	}

	Fixtures, err = testfixtures.New(
		testfixtures.Database(conn),
		testfixtures.Dialect("mysql"),
		testfixtures.Directory("../fixtures"),
		testfixtures.SkipResetSequences(),
	)

	if err != nil {
		panic(err)
	}
}
