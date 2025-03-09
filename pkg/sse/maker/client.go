package maker

import (
	"go-rest-api-boilerplate/api/sseutil"
	"go-rest-api-boilerplate/types"
)

type IClient interface {
	ReceiveMessage()
	GetSendChannel() chan *sseutil.SseRes
	Close()
	GetID() int64
	GetNickname() string
	UpdateResult(result *types.Results) error
	GetContestResults(ownerID int64) types.ContestItemResult
	IsConnected() bool
	SendMessage(message *sseutil.SseRes)
}
