package repo

import (
	user_repo "go-rest-api-boilerplate/internal/db/repo/user"
	"go-rest-api-boilerplate/internal/services/account"
)

type Store struct {
	AccountRepo user_repo.ReadWriter
	AccountAPI  *account.Service
}

func New(accountRepo user_repo.ReadWriter, accountAPI *account.Service) *Store {
	return &Store{
		AccountRepo: accountRepo,
		AccountAPI:  accountAPI,
	}
}
