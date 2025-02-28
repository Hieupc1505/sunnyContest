package account

import app "go-rest-api-boilerplate/internal"

const (
	MaxNicknameLen = 50
	MinNicknameLen = 4
)

type Nickname string

func (n Nickname) String() string {
	return string(n)
}

func NewNickname(nn string) (Nickname, error) {
	if len(nn) < MinNicknameLen || len(nn) > MaxNicknameLen {
		return "", app.ErrInvalidData
	}
	return Nickname(nn), nil
}
