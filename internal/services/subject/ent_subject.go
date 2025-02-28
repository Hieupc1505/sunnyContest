package subject

import (
	"github.com/jackc/pgx/v5/pgtype"
	db "go-rest-api-boilerplate/internal/db/sqlc"
)

type AddSubjectParams struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func NewSfSubject(id int64, userId int64, name string, description pgtype.Text, tags pgtype.Text) db.SfSubject {
	return db.SfSubject{
		ID:          id,
		UserID:      userId,
		Name:        name,
		Description: description,
		Tags:        tags,
	}
}

func NewAddSubjectParams(name Name, description string, tags Tags) *db.AddSubjectParams {
	des := pgtype.Text{
		String: description,
		Valid:  true,
	}
	tagJ := pgtype.Text{
		String: tags.String(),
		Valid:  true,
	}
	return &db.AddSubjectParams{
		Name:        name.String(),
		Description: des,
		Tags:        tagJ,
	}
}

func NewUpdateSubjectParams(s *db.AddSubjectParams, id int64) *db.UpdateSubjectParams {
	des := pgtype.Text{
		String: s.Description.String,
		Valid:  true,
	}
	tagJ := pgtype.Text{
		String: s.Tags.String,
		Valid:  true,
	}

	return &db.UpdateSubjectParams{
		ID:          id,
		Name:        s.Name,
		Description: des,
		Tags:        tagJ,
	}
}
