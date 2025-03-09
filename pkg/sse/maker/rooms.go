package maker

import (
	"go-rest-api-boilerplate/api/sseutil"
)

// RoomPublisher quản lý trạng thái publish/unpublish của một room
type RoomPublisher interface {
	Publisher(id int64)
	UnPublisher(id int64)
	IsPublished(id int64) bool
}

// RoomBroadcaster gửi tin nhắn đến một room cụ thể
type RoomBroadcaster interface {
	Broadcast(id int64, message *sseutil.SseRes)
}

// RoomManager quản lý danh sách rooms và truy xuất hub tương ứng
type RoomManager interface {
	IsExists(id int64) bool
	Add(roomID int64, hub IHub)
	Remove(roomID int64)
	GetHub(id int64) (IHub, error)
}

// IRoom kết hợp tất cả các interface trên
type IRoom interface {
	RoomPublisher
	RoomBroadcaster
	RoomManager
}
