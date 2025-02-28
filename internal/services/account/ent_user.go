package account

import (
	"github.com/jackc/pgx/v5/pgtype"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Role        int32     `json:"role"`
	Status      int32     `json:"status"`
	CreatedTime time.Time `json:"created_time"`
	Profile     db.AddProfileRow
}

func NewUser(username Username, password Password) *User {
	return &User{Username: username.String(), Password: password.String()}
}

func (u *User) ToParams() *db.AddParams {
	return &db.AddParams{
		Username: u.Username,
		Password: u.Password,
	}
}

func NewUpdateTokenParam(id int64, token string, exp time.Duration) (*db.UpdateUserTokenParams, error) {
	Token := pgtype.Text{
		String: token,
		Valid:  true,
	}

	Exp := pgtype.Timestamptz{
		Time:  time.Now().Add(exp),
		Valid: true,
	}

	return &db.UpdateUserTokenParams{
		ID:           id,
		Token:        Token,
		TokenExpired: Exp,
	}, nil

}

func ToUserInfo(u db.GetUserByUsernameRow) *UserInfo {
	return &UserInfo{
		ID:          u.ID,
		Username:    u.Username,
		Role:        u.Role,
		Status:      u.Status,
		CreatedTime: u.CreatedTime,
		Profile: db.AddProfileRow{
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
		},
	}
}
