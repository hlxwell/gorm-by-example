package models

import (
	"errors"
	"fmt"

	"github.com/hlxwell/gorm-by-example/db"
	"gorm.io/gorm"
)

type UserRole struct {
	gorm.Model
	UserID uint `gorm:"index"`
	RoleID uint `gorm:"index"`
	User   User
	Role   Role
}

// Validate uniqueness(userID, roleID)
func (r *UserRole) Validate() (err error) {
	err = db.Conn.Preload("User").Where("user_id = ? AND role_id = ?", r.UserID, r.RoleID).First(r).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return fmt.Errorf("duplicated Role %d for User %d", r.RoleID, r.UserID)
}

func (r *UserRole) BeforeCreate(tx *gorm.DB) (err error) {
	return r.Validate()
}

// After Create, Update and Delete will udpate counter
func (r *UserRole) AfterCreate(tx *gorm.DB) (err error) {
	if r.ID == 0 {
		return nil
	}

	tx.Preload("User").Where("user_id = ? AND role_id = ?", r.UserID, r.RoleID).First(r)
	return tx.Model(&User{}).Where("id = ?", r.UserID).Update("role_count", gorm.Expr("role_count + 1")).Error
}

func (r *UserRole) AfterDelete(tx *gorm.DB) (err error) {
	if r.ID == 0 {
		return nil
	}

	tx.Preload("User").Where("user_id = ? AND role_id = ?", r.UserID, r.RoleID).First(r)
	return tx.Model(&User{}).Where("id = ?", r.UserID).Update("role_count", gorm.Expr("role_count - 1")).Error
}
