package password

import (
	"regexp"
)

const (
	minPassLength = 6
	maxPassLength = 64
)

// Validate input to obey password requirement
// Passwords must be minimum 6 characters and maximum 64 characters,
// containing at least 1 capital characters AND 1 number AND 1 special (non alpha-numeric) characters.
func Validate(in string) bool {
	if len(in) < minPassLength {
		return false
	}

	if len(in) > maxPassLength {
		return false
	}

	if !hasCapitalCharacter(in) {
		return false
	}

	if !hasNumber(in) {
		return false
	}

	if !hasSpecialCharacter(in) {
		return false
	}

	return true
}

func hasCapitalCharacter(input string) bool {
	re := regexp.MustCompile(`[A-Z]`)
	return re.MatchString(input)
}

func hasNumber(input string) bool {
	re := regexp.MustCompile(`[0-9]`)
	return re.MatchString(input)
}

func hasSpecialCharacter(input string) bool {
	re := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	return re.MatchString(input)
}
