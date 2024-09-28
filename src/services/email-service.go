package services

import "regexp"

type EmailService struct{}

func InitEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) ValidateEmail(email string) bool {
	// Regular expression pattern for validating email addresses
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}
