package domain

import "gorm.io/gorm"

type UserClaim struct {
	gorm.Model
	UserId  uint
	ClaimId uint
}
