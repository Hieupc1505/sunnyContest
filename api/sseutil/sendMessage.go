package sseutil

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/response"
	"log"
)

var (
	ErrCodeSuccess = 0
	ErrCodeFail    = 1

	ErrInvalidContestID = 19
)

func SendMessage(ctx *gin.Context, errorCode *int, data any) error {
	if errorCode == nil {
		errorCode = &ErrCodeSuccess
	}
	event := response.Response{
		Error: *errorCode,
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

	_, err = ctx.Writer.Write(sseMessage)
	if err != nil {
		log.Printf("Failed to write message to client: %v\n", err)
		return err
	}
	ctx.Writer.Flush()
	return nil
}
