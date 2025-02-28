package subject

import (
	"encoding/json"
	app "go-rest-api-boilerplate/internal"
	"log/slog"
)

type Tags string

func (t Tags) String() string {
	return string(t)
}

func NewTags(t []string) (Tags, error) {
	tagJson, err := json.Marshal(t)
	if err != nil {
		slog.Info("Cannot marshal tags", "error", err)
		return "", app.ErrInternalServerError
	}
	return Tags(tagJson), nil
}
