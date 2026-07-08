package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

const (
	usernameMinLen = 3
	usernameMaxLen = 100
	fullNameMinLen = 3
	fullNameMaxLen = 100
	passwordMinLen = 6
	passwordMaxLen = 100
	emailMinLen    = 3
	emailMaxLen    = 320
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string, minLenght int, maxLenght int) error {
	l := len(value)

	if l < minLenght || l > maxLenght {
		return fmt.Errorf("must contain from %d-%d characters", minLenght, maxLenght)
	}

	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, usernameMinLen, usernameMaxLen); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return fmt.Errorf("must contain only lowercase letters, digits, or underscore")
	}

	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, fullNameMinLen, fullNameMaxLen); err != nil {
		return err
	}

	if !isValidFullName(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}

	return nil
}

func ValidatePassword(value string) error {
	if err := ValidateString(value, passwordMinLen, passwordMaxLen); err != nil {
		return err
	}

	return nil
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, emailMinLen, emailMaxLen); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not a valid email address")
	}

	return nil
}
