package helpers

import "golang.org/x/crypto/bcrypt"

type CryptoHelper struct {
}

func InitCryptoHelper() *CryptoHelper {
	return &CryptoHelper{}
}

func (h *CryptoHelper) Encrypt(s string) (string, error) {
	hashedString, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedString), nil
}

func (h *CryptoHelper) IsHashMatched(hash, s string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s)); err != nil {
		return false
	}
	return true
}
