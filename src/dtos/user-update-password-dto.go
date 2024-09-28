package dtos

type UserUpdatePasswordDto struct {
	UserId      uint
	NewPassword string
	OldPassword string
}
