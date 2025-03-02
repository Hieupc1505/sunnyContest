package repo

import (
	question_repo "go-rest-api-boilerplate/internal/db/repo/question"
	subject_repo "go-rest-api-boilerplate/internal/db/repo/subject"
	user_repo "go-rest-api-boilerplate/internal/db/repo/user"
	"go-rest-api-boilerplate/internal/services/account"
	"go-rest-api-boilerplate/internal/services/question"
	"go-rest-api-boilerplate/internal/services/subject"
)

type Store struct {
	AccountRepo     user_repo.ReadWriter
	AccountAPI      *account.Service
	SubjectRepo     subject_repo.ReadWriter
	SubjectService  *subject.Service
	QuestionRepo    question_repo.ReadWriter
	QuestionService *question.Service
}

func New(accountRepo user_repo.ReadWriter, accountAPI *account.Service, subRepo subject_repo.ReadWriter, subService *subject.Service, qr question_repo.ReadWriter, qs *question.Service) *Store {
	return &Store{
		AccountRepo:     accountRepo,
		AccountAPI:      accountAPI,
		SubjectRepo:     subRepo,
		SubjectService:  subService,
		QuestionRepo:    qr,
		QuestionService: qs,
	}
}
