package models

import (
	"github.com/hlxwell/gorm-by-example/db"
	"github.com/hlxwell/gorm-by-example/test"
)

func prepareTestDB() {
	db.InitTestDB()
	MigrateSchema()
	test.LoadFixtures()
}

// MigrateSchema to auto migrate all schema
// To avoid cycle import, put the migration here:
// model test -> db -> model
// model test -> test -> db, model
func MigrateSchema() {
	err := db.Conn.AutoMigrate(
		&User{},
		&Role{},
		&Feature{},
		&UserRole{},
	)
	if err != nil {
		panic(err)
	}

	err = db.Conn.SetupJoinTable(&User{}, "Roles", &UserRole{})
	if err != nil {
		panic(err)
	}
}
