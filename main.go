package main

import (
	"github.com/hlxwell/gorm-by-example/db"
	"github.com/hlxwell/gorm-by-example/models"
)

func main() {
	db.InitDB()
	models.MigrateSchema()
}
