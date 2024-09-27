package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `gorm:"unique"`
	Password     string
	EmailAddress string `gorm:"uniqueIndex"`
}
