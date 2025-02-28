package account

import "strings"

const (
	minUsernameLen = 2
	maxUsernameLen = 100
)

type Username string

func (e Username) String() string {
	return string(e)
}

func NewUsername(email string) (Username, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return "", ErrUsernameEmpty
	}
	if len(email) < minUsernameLen || len(email) > maxUsernameLen {
		return "", ErrUsernameInvalidLen
	}
	// add new rule
	return Username(email), nil
}
