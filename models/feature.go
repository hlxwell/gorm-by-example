package models

import "gorm.io/gorm"

type Feature struct {
	gorm.Model
	ID     uint   `gorm:"primaryKey"`
	Name   string `gorm:"unique"`
	RoleID uint
	Role   Role
}
