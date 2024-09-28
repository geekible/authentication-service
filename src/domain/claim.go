package domain

import "gorm.io/gorm"

type Claim struct {
	gorm.Model
	Claim string `gorm:"unique"`
}
