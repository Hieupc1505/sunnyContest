package sse

import (
	"go-rest-api-boilerplate/types"
	"sync"
)

type Room struct {
	State   types.SseStatus
	Host    chan *Client
	Clients map[int64]*Hub
	mu      sync.RWMutex
}

type State struct {
	Members map[int64]chan sse.SseStatus // Map of user ID to their message channel
	State   chan bool                    // A channel to hold any state information
}

var (
	instance IRoomManager
	once     sync.Once // Ensures Rooms map initialization occurs only once
)

type Subscriber struct {
	rdb *redis.Client
}

type Publisher struct {
	rdb *redis.Client
}

type IRoomMaker interface {
	Publisher(roomID int64)
	UnPublisher(roomID int64)
	UnSubscribe(roomID int64, userID int64)
	Subscribe(user int64, roomID int64, ch chan sse.SseStatus) interface{}
	CloseAfter(roomID int64, duration int)
}

type IRoomAction interface {
	IsRoomNotExist(roomID int64) bool
	GetState(roomID int64) chan bool
	ChangeState(roomID int64, status bool)
}

type IRoomMessage interface {
	NotifyAll(roomID int64, messages ...sse.SseStatus) error
}

type IRoomManager interface {
	IRoomMaker
	IRoomAction
	IRoomMessage
}

//	func init() {
//		once.Do(func() {
//			//TODO: Change the implementation to use the PubSub implementation
//			instance = NewManagerPubSub(share.GetRedis())
//		})
//	}
func NewRoomManager(client *redis.Client) IRoomManager {
	once.Do(func() {
		//TODO: Change the implementation to use the PubSub implementation
		instance = NewManagerPubSub(client)
		// instance = NewManager()
	})
	return instance
}
