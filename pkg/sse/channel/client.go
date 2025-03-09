package channel

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/response"
	"go-rest-api-boilerplate/api/sseutil"
	"go-rest-api-boilerplate/pkg/sse/maker"
	"go-rest-api-boilerplate/types"
	"log"
	"sync"
	"time"
)

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	ID       int64
	NickName string
	Results  *types.Results
	Ctx      *gin.Context
	Send     chan *sseutil.SseRes
	mu       sync.Mutex
	quit     chan struct{}
	once     sync.Once
}

func (c *Client) SendMessage(message *sseutil.SseRes) {
	c.Send <- message
}

func NewClient(ctx *gin.Context, ID int64, NickName string, bufferSize int) maker.IClient {
	return &Client{
		ID:       ID,
		Send:     make(chan *sseutil.SseRes, bufferSize),
		Ctx:      ctx,
		NickName: NickName,
		quit:     make(chan struct{}),
	}
}

func (c *Client) GetContestResults(ownerID int64) types.ContestItemResult {
	c.mu.Lock()
	defer c.mu.Unlock()
	var results types.ContestItemResult
	results.ID = c.ID
	results.Nickname = c.NickName
	results.Owner = ownerID == c.ID
	if c.Results != nil {
		results.Results = c.Results
	}

	return results
}

func (c *Client) IsConnected() bool {
	// Kiểm tra trạng thái kết nối của client
	// (ví dụ: kiểm tra kênh Send có đóng không)
	return c.Send != nil
}

func (c *Client) UpdateResult(result *types.Results) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Results = result
	return nil
}

// GetSendChannel trả về kênh Send của client.
func (c *Client) GetSendChannel() chan *sseutil.SseRes {
	return c.Send
}

// GetID returns userId
func (c *Client) GetID() int64 {
	return c.ID
}

// GetNickName returns userId
func (c *Client) GetNickname() string {
	return c.NickName
}

// Close đóng kênh Send của client.
func (c *Client) Close() {
	c.once.Do(func() {
		close(c.quit) // Đảm bảo chỉ đóng 1 lần
		fmt.Println("Channel closed")
	})
}

// writeMessage pumps messages from the hub to the client.
func (c *Client) writeMessage(data *sseutil.SseRes) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	statusCode := 0
	if v, ok := data.PkgData.(string); ok && v == string(types.CloseContest) {
		statusCode = 1
	}
	event := response.Response{
		Error: statusCode,
		Data:  data,
	}

	// Chuyển `data` thành JSON string
	message, err := json.Marshal(event)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return err
	}

	// Format SSE: "data: <json>\n\n"
	sseMessage := append([]byte("data: "), message...)
	sseMessage = append(sseMessage, []byte("\n\n")...)

	_, err = c.Ctx.Writer.Write(sseMessage)
	if err != nil {
		log.Printf("Failed to write message to client: %v\n", err)
		return err
	}
	c.Ctx.Writer.Flush()

	return nil
}

func (c *Client) ReceiveMessage() {
	defer func() {
		log.Println("Client ReceiveMessage disconnected")
	}()

	for {
		select {
		case message := <-c.Send:
			// Gửi tin nhắn đến client.
			if err := c.writeMessage(message); err != nil {
				log.Println("Error writing to client:", err)
				return
			}
		case <-c.quit:
			return
		}
	}
}
