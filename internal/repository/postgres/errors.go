package postgres

import (
	"strings"
)

var (
	ErrEventNotFound          = "event not found"
	ErrEventNameAlreadyExists = "event name already exists"
	ErrClubNameAlreadyExists  = "club name already exists"
	ErrClubNotFound           = "club not found"
)

func IsDuplicateKeyError(err error) (bool, string) {
	if strings.Contains(err.Error(), "duplicate key") {
		// Example: duplicate key value violates unique constraint "users_username_key"
		parts := strings.Split(err.Error(), `"`)
		if len(parts) == 3 {
			return true, strings.Split(parts[1], "_")[1]
		}
	}

	return false, ""
}
