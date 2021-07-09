package models

import (
	"fmt"

	"github.com/hlxwell/gorm-by-example/db"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string  `gorm:"unique"`
	Roles      []*Role `gorm:"many2many:user_roles"`
	RoleCount  uint
	UserRoles  []*UserRole
	Attributes datatypes.JSON
}

func WithRole(role string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Joins("INNER JOIN user_roles ON user_roles.user_id = users.id").
			Joins("INNER JOIN roles ON user_roles.role_id = roles.id").
			Where("roles.name = ?", role)
	}
}

func (u *User) AddRole(role Role) error {
	if role.ID == 0 {
		return fmt.Errorf("missing Role.ID when add role")
	}

	return db.Conn.Model(&u).Association("UserRoles").Append(&UserRole{RoleID: role.ID})
}

func (u *User) DeleteRole(role Role) error {
	if role.ID == 0 {
		return fmt.Errorf("missing Role.ID when delete role")
	}

	var userRole UserRole
	db.Conn.Where("user_id = ? AND role_id = ?", u.ID, role.ID).First(&userRole)
	return db.Conn.Debug().Delete(&userRole).Error
}
