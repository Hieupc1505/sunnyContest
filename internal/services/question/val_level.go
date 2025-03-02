package question

import (
	app "go-rest-api-boilerplate/internal"
	db "go-rest-api-boilerplate/internal/db/sqlc"
)

type Level string

func (l *Level) String() string {
	return string(*l)
}

var validLevelQuestions = map[db.LevelQuestion]struct{}{
	db.LevelQuestionEASY:   {},
	db.LevelQuestionMEDIUM: {},
	db.LevelQuestionHARD:   {},
}

func NewLevel(l string) (Level, error) {
	_, exists := validLevelQuestions[db.LevelQuestion(l)]
	if !exists {
		return "", app.ErrInvalidData
	}
	return Level(l), nil
}
