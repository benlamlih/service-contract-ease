package domain

import (
	"crypto/rand"
	"fmt"
	"strings"
)

type CreateUserParams struct {
	Username  string
	Email     string
	Password  string
	FirstName string
	LastName  string
}

func GenerateUsername(firstName, lastName string) string {
	first := strings.ToLower(strings.ReplaceAll(firstName, " ", ""))
	last := strings.ToLower(strings.ReplaceAll(lastName, " ", ""))

	suffix := make([]byte, 3)
	if _, err := rand.Read(suffix); err != nil {
		suffix = []byte{1, 2, 3}
	}

	return fmt.Sprintf("%s.%s%05d", first, last, int(suffix[0])<<16|int(suffix[1])<<8|int(suffix[2]))
}
