package util

import (
	"errors"
	"fmt"
	"regexp"
)

// GetRandomAvatar returns a random avatar URL from the randomuser.me API
func GetRandomAvatar(index int) string {
	return fmt.Sprintf("https://randomuser.me/api/portraits/lego/%d.jpg", index)
}

// ValidateUserName checks if the username is 4-16 characters long and contains only alphanumeric characters.
func ValidateUserName(username string) error {
	if len(username) < 4 || len(username) > 16 {
		return errors.New("username must be between 4 and 16 characters")
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, username)
	if !matched {
		return errors.New("username can only contain alphanumeric characters")
	}

	return nil
}
