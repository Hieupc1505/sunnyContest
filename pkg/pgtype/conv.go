package pgtype_conv

import "github.com/jackc/pgx/v5/pgtype"

func NewString(s string) pgtype.Text {
	return pgtype.Text{
		Valid:  true,
		String: s,
	}
}
