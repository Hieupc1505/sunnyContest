package subject

import (
	"github.com/gin-gonic/gin"
	app "go-rest-api-boilerplate/internal"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/internal/services"
	"go-rest-api-boilerplate/pkg/errsx"
	"log/slog"
)

func makeSubject(name string, des string, tags []string) (*db.AddSubjectParams, error) {
	var errs errsx.Map
	n, err := NewName(name)
	if err != nil {
		errs.Set("Subject_name", app.ErrInvalidData)
	}

	t, err := NewTags(tags)
	if err != nil {
		errs.Set("Subject_tags", app.ErrInvalidData)
	}

	if len(errs) > 0 {
		return nil, errs.ToError()
	}

	sub := NewAddSubjectParams(n, des, t)

	return sub, nil
}

func (s *Service) AddSubject(ctx *gin.Context, userId int64, name, description string, tags []string) (*db.AddSubjectRow, error) {
	sub, err := makeSubject(name, description, tags)
	if err != nil {
		return nil, services.NewError(err, nil)
	}

	sub.UserID = userId
	sub.State = StateDefault

	// Lưu user vào db
	sb, err := s.repo.AddSubject(ctx, *sub)
	if err != nil {
		slog.Info("Error adding subject", "error", err)
		return nil, services.NewError(err, app.ErrInternalServerError)
	}

	return &sb, nil
}

func (s *Service) Update(ctx *gin.Context, subjectID int64, name, description string, tags []string) (*db.UpdateSubjectRow, error) {
	sub, err := makeSubject(name, description, tags)
	if err != nil {
		return nil, services.NewError(err, nil)
	}
	updateParams := NewUpdateSubjectParams(sub, subjectID)
	subUpdated, err := s.repo.UpdateSubject(ctx, *updateParams)
	if err != nil {
		slog.Info("Error updating subject", "error", err)
		return nil, services.NewError(err, app.ErrInternalServerError)
	}
	return &subUpdated, nil

}
