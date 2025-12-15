package auth

import (
	"errors"
	"regexp"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9]{3,15}$`)

	passwordLenRegex     = regexp.MustCompile(`^\S{8,}$`)
	passwordDigitRegex   = regexp.MustCompile(`[0-9]`)
	passwordSpecialRegex = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?~]`)
)

var (
	ErrInvalidUsername = errors.New("username must be 3-15 characters and contain only letters and numbers")
	ErrInvalidPassword = errors.New("password must be at least 8 characters and include a number and a special character")
)

func ValidateUsername(username string) error {
	if !usernameRegex.MatchString(username) {
		return ErrInvalidUsername
	}
	return nil
}

func ValidatePassword(password string) error {
	if !(passwordLenRegex.MatchString(password) && passwordDigitRegex.MatchString(password) && passwordSpecialRegex.MatchString(password)) {
		return ErrInvalidPassword
	}
	return nil
}
