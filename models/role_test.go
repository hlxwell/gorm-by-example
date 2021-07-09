package models

import (
	"testing"

	"github.com/hlxwell/gorm-by-example/db"
	"github.com/hlxwell/gorm-by-example/test"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/clause"
)

func TestRoleSpec(t *testing.T) {
	prepareTestDB()

	Convey("With Fixtures", t, func() {
		test.LoadFixtures()

		Convey("Should be able to create Role", func() {
			result := db.Conn.Create(&Role{
				Name: "admin",
			})
			So(result.Error, ShouldBeNil)
		})

		Convey("Should be able to Presence Validation for Role.Name before save", func() {
			var role Role
			db.Conn.First(&role)
			result := db.Conn.Create(&Role{
				Name: role.Name,
			})
			So(result.Error, ShouldNotBeNil)
		})

		Convey("Should be able to Check Role.Users", func() {
			var role Role
			db.Conn.Preload("Users").First(&role)
			So(len(role.Users), ShouldBeGreaterThan, 0)
		})

		Convey("Should be able to find normal role by scope", func() {
			var role Role
			So(db.Conn.Scopes(IsNormal).First(&role).Error, ShouldBeNil)
			So(role.ID, ShouldBeGreaterThan, 0)
		})
	})
}

// after role delete, the user.RoleCount should be affected.
func TestDeleteRole(t *testing.T) {
	prepareTestDB()

	var user User
	db.Conn.Preload(clause.Associations).First(&user)
	for _, role := range user.Roles {
		assert.NoError(t, role.Destroy())
	}

	var newUser User
	db.Conn.First(&newUser)
	assert.Equal(t,
		db.Conn.Model(&newUser).Association("Roles").Count(),
		int64(newUser.RoleCount),
	)
}
