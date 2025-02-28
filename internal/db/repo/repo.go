package repo

import (
	subject_repo "go-rest-api-boilerplate/internal/db/repo/subject"
	user_repo "go-rest-api-boilerplate/internal/db/repo/user"
	"go-rest-api-boilerplate/internal/services/account"
	"go-rest-api-boilerplate/internal/services/subject"
)

type Store struct {
	AccountRepo    user_repo.ReadWriter
	AccountAPI     *account.Service
	SubjectRepo    subject_repo.ReadWriter
	SubjectService *subject.Service
}

func New(accountRepo user_repo.ReadWriter, accountAPI *account.Service, subRepo subject_repo.ReadWriter, subService *subject.Service) *Store {
	return &Store{
		AccountRepo:    accountRepo,
		AccountAPI:     accountAPI,
		SubjectRepo:    subRepo,
		SubjectService: subService,
	}
}
