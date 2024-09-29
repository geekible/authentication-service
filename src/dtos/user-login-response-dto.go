package dtos

import "time"

type UserLoginResponseDto struct {
	Username     string
	EmailAddress string
	UserClaims   []string
	Exp          time.Time
}
