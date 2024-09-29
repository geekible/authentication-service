package dtos

type UserLoginResponseDto struct {
	Username     string
	EmailAddress string
	UserClaims   []string
}
