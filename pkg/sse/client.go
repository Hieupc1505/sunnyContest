package sse

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/sseutil"
	"go-rest-api-boilerplate/types"
	"log"
	"log/slog"
	"time"
)

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	Hub  *Hub
	Sse  chan []byte
	Ctx  *gin.Context
	Send chan []byte
}

func (h *Hub) Control() {
	for {
		select {
		case <-h.done:
			h.EndContest()

			time.Sleep(1 * time.Minute)
			h.close <- struct{}{}
		case <-h.close:
			done := make(chan struct{}) // Tạo channel để đồng bộ hóa
			go func() {
				endMessage := sseutil.NewSseRes(types.CloseContest, gin.H{
					"id": h.ID,
				})
				h.Broadcast(endMessage)
				close(done) // Đóng channel khi Broadcast xong
			}()

			<-done // Đợi Broadcast hoàn thành
			time.Sleep(10 * time.Second)
			close(h.done)
			h.Close()
			return
		case <-h.quit:
			return
		}
	}
	fmt.Println("End control")

}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.registerChan:
			// Đăng ký client mới
			h.mu.Lock()
			h.Clients[client.GetID()] = client
			h.mu.Unlock()

		case client := <-h.unregisterChan:
			// Hủy đăng ký client
			h.mu.Lock()
			id := client.GetID()
			if _, ok := h.Clients[id]; ok {
				delete(h.Clients, id)
				client.Close() // Đóng kết nối của client
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for id, client := range h.Clients {
				select {
				case client.GetSendChannel() <- message:
					// Tin nhắn được gửi thành công
				default:
					// Kênh của client đã đầy, bỏ qua client này
					// Không đóng kết nối, chỉ ghi log để theo dõi
					slog.Warn("Client channel is full, skipping message", "clientID", id)
				}
			}
			h.mu.RUnlock()
		case <-h.quit:
			fmt.Println("end run goroutine")
			return
		}
	}
	fmt.Println("End run")
}

func (h *Hub) StartTimer(timeExam int32) {
	// Tạo một bộ đếm thời gian mỗi 15 giây
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop() // Dừng bộ đếm khi hàm kết thúc

	// Tính thời gian kết thúc (ví dụ: 5 phút sau)
	duration := time.Duration(timeExam) * time.Minute
	startTime := time.Now().Add(duration)
	endTime := startTime.Add(duration)

	// Gửi tin nhắn bắt đầu cuộc thi
	data := sseutil.NewSseRes(types.StartContest, gin.H{
		"id":         h.ID,
		"time_start": endTime.Unix(),
	})
	h.Broadcast(data)

	h.timeStart = endTime

	// Tính thời gian còn lại
	remainingTime := time.Duration(endTime.Unix()-startTime.Unix()) * time.Second

	// Tạo một timer chính xác cho thời điểm kết thúc
	endTimer := time.After(remainingTime)

	for {
		select {
		case <-h.done:
			fmt.Println("end end timer goroutine")
			return
		case <-ticker.C:
			// Tính thời gian còn lại
			remainingTime := endTime.Unix() - time.Now().Unix()

			// Gửi tin nhắn định kỳ mỗi 15 giây
			message := sseutil.NewSseRes(types.ContestCountDown, gin.H{
				"id":        h.ID,
				"countdown": remainingTime,
			})
			h.Broadcast(message)

		case <-endTimer:
			h.Stop()
			// Dừng bộ đếm
			return
		}
	}
	slog.Info("End Start timer goroutine")
}

func (h *Hub) EndContest() {

	// Gửi tin nhắn kết thúc cuộc thi
	endMessage := sseutil.NewSseRes(types.EndContest, gin.H{
		"id": h.ID,
	})
	h.Broadcast(endMessage)

	err := h.handler.ContestRepo.StopContest(context.Background(), h.ID)
	if err != nil {
		slog.Info("Update Stop contest fail", "error", err)
	}

}

func (h *Hub) Stop() {
	h.done <- struct{}{}
}
