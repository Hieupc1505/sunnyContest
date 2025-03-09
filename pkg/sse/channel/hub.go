package channel

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/handler"
	"go-rest-api-boilerplate/api/sseutil"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/internal/services/question"
	"go-rest-api-boilerplate/pkg/sse/maker"
	"go-rest-api-boilerplate/types"
	"log/slog"
	"sync"
	"time"
)

type Hub struct {
	ID             int64
	Clients        map[int64]maker.IClient
	Published      bool
	broadcast      chan *sseutil.SseRes
	registerChan   chan maker.IClient
	unregisterChan chan maker.IClient
	Questions      []types.Question //int64 -> id của question, với giá trị là index của answer chính xác
	timeStart      time.Time
	mu             sync.RWMutex
	handler        *handler.Handler
	once           sync.Once
	wg             sync.WaitGroup

	done      chan struct{}
	stopTimer chan struct{} // Channel để nhận tín hiệu dừng
	close     chan struct{}
	quit      chan struct{}
}

func NewHub(handler *handler.Handler, contestID int64) maker.IHub {
	return &Hub{
		ID:             contestID,
		Clients:        make(map[int64]maker.IClient),
		broadcast:      make(chan *sseutil.SseRes, 10),
		registerChan:   make(chan maker.IClient),
		unregisterChan: make(chan maker.IClient),
		done:           make(chan struct{}, 1),
		close:          make(chan struct{}, 1),
		quit:           make(chan struct{}),
		handler:        handler,
	}
}

func (h *Hub) Done() <-chan struct{} {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.done
}

func (h *Hub) Quit() <-chan struct{} {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.quit
}

func (h *Hub) UserSubmit(id int64, results *types.Results) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clients, exists := h.Clients[id]
	if !exists {
		return fmt.Errorf("cannot find client from hub %d", id)
	}

	if err := clients.UpdateResult(results); err != nil {
		slog.Error("Error updating user's result", "error", err)
		return err
	}

	return nil

}

func (h *Hub) SetQuestions(questions []db.SfQuestion) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	data := make([]types.Question, len(questions))
	for i, q := range questions {
		data[i] = *question.NewQuestionFromSfQuestion(q) // Gán trực tiếp vào vị trí i
	}
	h.Questions = data
}
func (h *Hub) GetQuestions() []types.Question {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.Questions
}

// GetUsers updates the
func (h *Hub) GetUsers(ownerID int64) []types.ContestItemResult {
	h.mu.RLock()
	defer h.mu.RUnlock()
	var users []types.ContestItemResult
	for _, clients := range h.Clients {
		data := clients.GetContestResults(ownerID)
		users = append(users, data)
	}
	return users
}

// Publish : Đánh dấu Hub là Published
func (h *Hub) Publish() {
	h.mu.Lock()
	h.Published = true
	h.mu.Unlock()
}

// UnPublish : Hủy publish của Hub
func (h *Hub) UnPublish() {
	h.mu.Lock()
	h.Published = false
	h.mu.Unlock()
}

// IsPublished : check whether the Hub is published
func (h *Hub) IsPublished() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.Published
}

func (h *Hub) Register(client maker.IClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.registerChan <- client
}

func (h *Hub) UnRegister(client maker.IClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.unregisterChan <- client
}

func (h *Hub) Broadcast(message *sseutil.SseRes) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	h.broadcast <- message

}
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.registerChan:
			h.mu.Lock()
			h.Clients[client.GetID()] = client
			h.mu.Unlock()

		case client := <-h.unregisterChan:
			h.mu.Lock()
			id := client.GetID()
			if _, ok := h.Clients[id]; ok {
				//delete(h.Clients, id)
				client.Close()
			}
			h.mu.Unlock()

		case message, ok := <-h.broadcast:
			h.mu.RLock()
			if !ok {
				return
			}
			for _, client := range h.Clients {
				select {
				case client.GetSendChannel() <- message:
				default:
					//slog.Warn("Client channel is disconnect, skipping message", "clientID", id)
					continue
				}
			}
			h.mu.RUnlock()
		case <-h.done:
			h.mu.Lock()
			h.Stop()
			h.mu.Unlock()
			return
		case <-h.quit:
			h.CleanupClients()
			h.CleanupHub()
			return
		}
	}
}

func (h *Hub) StartTimer(startTime time.Time, timeExam int32) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	duration := time.Duration(timeExam) * time.Minute
	//startTime := time.Now()
	endTime := startTime.Add(duration)

	h.timeStart = endTime

	// Gửi thông báo bắt đầu contest
	data := sseutil.NewSseRes(types.StartContest, gin.H{
		"id":         h.ID,
		"time_start": endTime.UnixMilli(),
	})
	h.Broadcast(data)

	endTimer := time.NewTimer(duration)
	defer endTimer.Stop()
	h.stopTimer = make(chan struct{})

	for {
		select {
		case <-h.quit:
			return
		case <-ticker.C:
			remainingTime := endTime.Unix() - time.Now().Unix()
			if remainingTime <= 0 {
				continue
			}

			message := sseutil.NewSseRes(types.ContestCountDown, gin.H{
				"id":        h.ID,
				"countdown": remainingTime,
			})
			h.Broadcast(message)

		case <-endTimer.C:
			// Gửi thông báo EndContest khi hết thời gian thi
			h.EndContest()
			h.Close()
			return
		case <-h.stopTimer:
			// Gửi thông báo EndContest khi hết thời gian thi
			h.EndContest()
			h.Close()
			return
		}
	}
}

func (h *Hub) StopTimer() {
	close(h.stopTimer)
}

func (h *Hub) Close() {
	// Đợi 1 phút trước khi gửi thông báo ClosedContest
	timeWaitToSeeResult := 10 * time.Second
	timeHandleSendClosedMessage := 5 * time.Second
	time.Sleep(timeWaitToSeeResult)

	// Gửi thông báo ClosedContest
	h.Broadcast(sseutil.NewSseRes(types.CloseContest, gin.H{
		"id": h.ID,
	}))
	time.Sleep(timeHandleSendClosedMessage)
	h.done <- struct{}{}
}

func (h *Hub) EndContest() {

	endMessage := sseutil.NewSseRes(types.EndContest, gin.H{
		"id": h.ID,
	})
	h.Broadcast(endMessage)

	if err := h.handler.ContestRepo.StopContest(context.Background(), h.ID); err != nil {
		slog.Warn("Failed to update contest status", "error", err)
	}

}

func (h *Hub) Stop() {
	h.once.Do(func() {
		close(h.quit)
	})
}

func (h *Hub) CleanupHub() {
	h.mu.Lock()
	defer h.mu.Unlock()

	close(h.broadcast)
	close(h.registerChan)
	close(h.unregisterChan)
	h.Stop()
}

func (h *Hub) CleanupClients() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, client := range h.Clients {
		//delete(h.Clients, client.GetID())
		client.Close()
	}
	h.Clients = make(map[int64]maker.IClient) // Xóa toàn bộ client
}
