package helpers

import (
	"fmt"
	"unicode"
)

type PasswordHelper struct {
	password          string
	minCharCount      int
	minLowerCaseCount int
	minUpperCaseCount int
	minSymbolsCount   int
}

func InitPasswordHelper(password string) *PasswordHelper {
	return &PasswordHelper{
		password:          password,
		minCharCount:      8,
		minLowerCaseCount: 1,
		minUpperCaseCount: 1,
		minSymbolsCount:   1,
	}
}

func (h *PasswordHelper) ValidateComplexity() error {
	lowerCaseLetterCount := 0
	upperCaseLetterCount := 0
	symbolsCount := 0

	for _, char := range h.password {
		if unicode.IsLower(char) {
			lowerCaseLetterCount++
		}
		if unicode.IsUpper(char) {
			upperCaseLetterCount++
		}
		if unicode.IsSymbol(char) {
			symbolsCount++
		}
	}

	symbolsCount++

	if lowerCaseLetterCount < h.minLowerCaseCount ||
		upperCaseLetterCount < h.minUpperCaseCount ||
		symbolsCount < h.minSymbolsCount ||
		len(h.password) < h.minCharCount {
		return fmt.Errorf("password must be at least %d character long and contain a mixture of upper and lowercase letter and at least %d symbols",
			h.minCharCount,
			h.minSymbolsCount)
	}

	return nil
}
