package channel

import (
	"fmt"
	"go-rest-api-boilerplate/api/sseutil"
	app "go-rest-api-boilerplate/internal"
	"go-rest-api-boilerplate/pkg/sse/maker"
	"sync"
)

type Rooms struct {
	rooms map[int64]maker.IHub
	mu    sync.RWMutex // guards
	once  sync.Once
}

// NewRooms khởi tạo Rooms mới
func NewRooms() *Rooms {
	return &Rooms{
		rooms: make(map[int64]maker.IHub),
	}
}

// IsExists returns true if Room is
func (r *Rooms) IsExists(id int64) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.rooms[id]
	return exists
}

// Publisher Triển khai RoomPublisher
func (r *Rooms) Publisher(id int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if hub, exists := r.rooms[id]; exists {
		hub.Publish()
	}
}

// UnPublisher khai
func (r *Rooms) UnPublisher(id int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if hub, exists := r.rooms[id]; exists {
		hub.UnPublish()
	}
}

func (r *Rooms) IsPublished(id int64) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	hub, exists := r.rooms[id]
	return exists && hub.IsPublished()
}

// Add Triển khai RoomManager
func (r *Rooms) Add(roomID int64, hub maker.IHub) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rooms[roomID] = hub
}

// Remove Triển khai RoomRemover
func (r *Rooms) Remove(id int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.once.Do(func() {
		fmt.Println("removing room::::", id)
		delete(r.rooms, id)
	})
}

func (r *Rooms) GetHub(id int64) (maker.IHub, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if hub, exists := r.rooms[id]; exists {
		return hub, nil
	}
	return nil, app.ErrContestNotFound
}

// Broadcast Triển khai RoomBroadcaster
func (r *Rooms) Broadcast(id int64, message *sseutil.SseRes) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	hub, exists := r.rooms[id]
	if exists && hub.IsPublished() {
		hub.Broadcast(message)
	}
}
