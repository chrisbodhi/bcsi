package utils

import (
	"errors"
	"regexp"
	"strings"
)

// ValidateSet ensures we have a key string, an equal sign, and then anything else, really.
func ValidateSet(input string) error {
	parts := strings.Split(input, "=")

	if len(parts) != 2 {
		return errors.New("usage: `set KEY=VALUE`")
	}

	key, value := parts[0], parts[1]

	if key == "" || value == "" {
		return errors.New("usage: `set KEY=VALUE`")
	}

	wordPattern := "^[a-zA-Z]+$"
	regex := regexp.MustCompile(wordPattern)

	if !regex.MatchString(key) {
		return errors.New("key must be a single word made of letters")
	}

	return nil
}

func WithSpace(inputs []string) string {
	return strings.Join(inputs, " ")
}
