package account

import db "go-rest-api-boilerplate/internal/db/sqlc"

type RegisterNickNameParams struct {
	NickName string `json:"nickname"`
	Type     int32  `json:"type"`
}

func NewProfile(userID int64, nn Nickname, avt Avatar) *db.AddProfileParams {
	return &db.AddProfileParams{
		UserID:   userID,
		Nickname: nn.String(),
		Avatar:   avt.String(),
	}
}
