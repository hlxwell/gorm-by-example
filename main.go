package main

import (
	"net/http"

	"github.com/hlxwell/gorm-by-example/auditable"
	"github.com/hlxwell/gorm-by-example/db"
	"github.com/hlxwell/gorm-by-example/models"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()
	db.InitDB()
	models.MigrateSchema()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("current_user_id", "12344321")
			return next(c)
		}
	})
	e.Use(auditable.GormInjector(db.Conn))
	e.GET("/hello", func(c echo.Context) error {
		conn := c.Get(auditable.GormDBKey).(*gorm.DB)
		conn.Create(&models.User{Name: "hello-hlxwell"})
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
