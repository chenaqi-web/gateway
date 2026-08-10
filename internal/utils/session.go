package utils

import (
	"errors"
	"regexp"
	"strings"
)

var sessionPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{16,80}$`)

func ValidateSession(sessionID string) error {
	if !sessionPattern.MatchString(strings.TrimSpace(sessionID)) {
		return errors.New("invalid upload session")
	}
	return nil
}
