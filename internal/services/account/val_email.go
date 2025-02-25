package account

import "strings"

const (
	minEmailLen = 2
	maxEmailLen = 100
)

type Email string

func (e Email) String() string {
	return string(e)
}

func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return "", ErrEmailEmpty
	}
	if len(email) < minEmailLen || len(email) > maxEmailLen {
		return "", ErrEmailInvalidLen
	}
	// add new rule
	return Email(email), nil
}
