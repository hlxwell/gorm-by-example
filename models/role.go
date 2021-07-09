package models

import (
	"github.com/hlxwell/gorm-by-example/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Role struct {
	gorm.Model
	Name      string  `gorm:"unique"`
	Users     []*User `gorm:"many2many:user_roles"`
	Features  []*Feature
	UserRoles []*UserRole
}

func IsNormal(db *gorm.DB) *gorm.DB {
	return db.Where("Name = ?", "normal")
}

// Destroy = delete self + delete relations
func (role *Role) Destroy() (err error) {
	db.Conn.Preload(clause.Associations).First(&role)
	return db.Conn.Transaction(func(tx *gorm.DB) error {
		if err = tx.Delete(&role.UserRoles).Error; err != nil {
			return err
		}

		if err = tx.Select(clause.Associations).Delete(&role).Error; err != nil {
			return err
		}

		return nil
	})
}
