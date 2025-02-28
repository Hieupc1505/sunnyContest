package subject

import (
	app "go-rest-api-boilerplate/internal"
)

const (
	MinSubLen = 4
)

type Name string

func (n Name) String() string {
	return string(n)
}

func NewName(s string) (Name, error) {
	if len(s) < 4 {
		return "", app.ErrInvalidData
	}
	return Name(s), nil
}
