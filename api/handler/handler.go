package handler

import (
	"github.com/gin-gonic/gin"
	contest_repo "go-rest-api-boilerplate/internal/db/repo/contest"
	question_repo "go-rest-api-boilerplate/internal/db/repo/question"
	subject_repo "go-rest-api-boilerplate/internal/db/repo/subject"
	user_repo "go-rest-api-boilerplate/internal/db/repo/user"
	user_contest_repo "go-rest-api-boilerplate/internal/db/repo/user_contest"
	"go-rest-api-boilerplate/internal/services/account"
	"go-rest-api-boilerplate/internal/services/contest"
	"go-rest-api-boilerplate/internal/services/question"
	"go-rest-api-boilerplate/internal/services/subject"
	"go-rest-api-boilerplate/internal/services/user_contest"
	"go-rest-api-boilerplate/pkg/imgUploader"
	"go-rest-api-boilerplate/pkg/sse/maker"
	"go-rest-api-boilerplate/pkg/token"
	"sync"
)

type Handler struct {
	App   *gin.Engine
	Tasks *sync.WaitGroup

	UserRepo       user_repo.ReadWriter
	AccountService *account.Service

	SubjectRepo    subject_repo.ReadWriter
	SubjectService *subject.Service

	QuestionRepo    question_repo.ReadWriter
	QuestionService *question.Service

	ContestRepo    contest_repo.ReadWriter
	ContestService *contest.Service

	UserContestRepo    user_contest_repo.ReadWriter
	UserContestService *user_contest.Service

	Rooms maker.IRoom

	Token       token.Maker
	ImgUploader imgUploader.IUploadImage
}
