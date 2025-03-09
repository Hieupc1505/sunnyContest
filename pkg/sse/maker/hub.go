package maker

import (
	"go-rest-api-boilerplate/api/sseutil"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/types"
	"time"
)

// ClientManager quản lý việc đăng ký / hủy đăng ký client
type ClientManager interface {
	Register(client IClient)
	UnRegister(client IClient)
	CleanupClients()
}

// Publisher quản lý trạng thái publish/unpublish
type Publisher interface {
	Publish()
	UnPublish()
	IsPublished() bool
}

// Broadcaster gửi tin nhắn đến tất cả client
type Broadcaster interface {
	Broadcast(message *sseutil.SseRes)
}

// HubRunner điều khiển vòng đời của Hub
type HubRunner interface {
	Run()
	Stop()
	Done() <-chan struct{}
	Quit() <-chan struct{}
	Close()
	//Control()
	//Leave()
	StopTimer()
	CleanupHub()
}

type ContestAction interface {
	GetUsers(ownerContestID int64) []types.ContestItemResult
	StartTimer(startTime time.Time, endTime int32)
	EndContest()
	GetQuestions() []types.Question
	SetQuestions(questions []db.SfQuestion)
	UserSubmit(userID int64, results *types.Results) error
}

// IHub kết hợp các interface trên
type IHub interface {
	ClientManager
	Publisher
	Broadcaster
	HubRunner
	ContestAction
}
