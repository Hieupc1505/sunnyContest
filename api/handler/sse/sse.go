package sse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/handler"
	"go-rest-api-boilerplate/api/response"
	"go-rest-api-boilerplate/api/sseutil"
	app "go-rest-api-boilerplate/internal"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/internal/services/contest"
	"go-rest-api-boilerplate/pkg/sse/channel"
	"go-rest-api-boilerplate/pkg/sse/maker"
	"go-rest-api-boilerplate/pkg/token"
	"go-rest-api-boilerplate/types"
	"log/slog"
	"strconv"
	"time"
)

const (
	bufferSize = 5
)

func Routes(h *handler.Handler) {
	sseGroup := h.App.Group("/api/v1/sse")
	sseGroup.GET("", Connect(h))
}

func setupSSEHeaders(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
}

func parseContestIDAndToken(ctx *gin.Context, h *handler.Handler) (int64, *token.Payload, error) {
	contestID := ctx.Query("rid")
	contestInt, err := strconv.ParseInt(contestID, 10, 64)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid contest id query")
	}

	token := ctx.Query("token")
	claims, err := h.Token.VerifyToken(token)
	if err != nil {
		fmt.Println("Verify token error in sse:::", err)
		slog.ErrorContext(ctx, "Error verifying token", "error", err)
		return 0, nil, app.ErrPermissionDenied
	}

	return contestInt, claims, nil
}

func getContestAndUserInfo(ctx *gin.Context, h *handler.Handler, contestInt int64, userID int64) (*db.SfContest, *db.GetUserByIDRow, error) {
	contestInfo, err := h.ContestRepo.GetContestLiveByID(ctx, contestInt)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting contest info", "error", err)
		return nil, nil, app.ErrContestNotFound
	}

	userInfo, err := h.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, nil, app.ErrUserNotFound
	}
	data := contest.ToSfContest(&contestInfo)
	return data, &userInfo, nil
}

func registerClientAndSendMessages(hub maker.IHub, client maker.IClient, contestInfo *db.SfContest) {
	hub.Register(client)
	go client.ReceiveMessage()

	//client.Send <- sseutil.NewSseRes(types.ContestInfo, contestInfo)
	message := sseutil.NewSseRes(types.ContestInfo, contestInfo)
	client.SendMessage(message)
	time.Sleep(100 * time.Millisecond)
	hub.Broadcast(sseutil.NewSseRes(types.UserJoin, hub.GetUsers(contestInfo.UserID)))
}

func handleDisconnect(hub maker.IHub, client maker.IClient) {
	hub.UnRegister(client)
}

// handleError xử lý lỗi và gửi phản hồi SSE
func handleError(ctx *gin.Context, err error) {
	rsp := response.ErrorResponse{Msg: app.ErrInvalidContestGameID.Error()}
	sseutil.SendMessage(ctx, &sseutil.ErrInvalidContestID, rsp)
}

func Connect(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Thiết lập header SSE
		setupSSEHeaders(ctx)

		// Xử lý contest ID và token
		contestInt, claims, err := parseContestIDAndToken(ctx, h)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		// Lấy thông tin contest và user
		contestInfo, userInfo, err := getContestAndUserInfo(ctx, h, contestInt, claims.UserID)
		if err != nil {
			handleError(ctx, err)
			return
		}

		if userInfo.Role == types.TeacherRole && contestInfo.UserID != userInfo.ID {
			response.Error(ctx, app.ErrPermissionDenied)
			return
		}

		// Lấy hub tương ứng với contest
		hub, err := h.Rooms.GetHub(contestInt)

		if err != nil {
			if !errors.Is(err, app.ErrContestNotFound) && contestInfo.UserID != userInfo.ID {
				handleError(ctx, err)
				return
			}
			hub = channel.NewHub(h, contestInt)
			hub.Publish()
			var data []db.SfQuestion
			if err := json.Unmarshal([]byte(contestInfo.Questions), &data); err != nil {
				handleError(ctx, err)
				return
			}
			hub.SetQuestions(data)
			h.Rooms.Add(contestInt, hub)
		}

		if contestInfo.UserID == userInfo.ID {
			go hub.Run()
		}

		if contestInfo.UserID != userInfo.ID {
			// Kiểm tra xem người dùng đã tham gia contest chưa
			userContest, err := h.UserContestRepo.GetUserContest(ctx, db.GetUserContestParams{
				ContestID: contestInt,
				UserID:    userInfo.ID,
			})

			if userContest.Result != nil || len(userContest.Result) != 0 {
				data := response.ErrorResponse{
					Msg: app.ErrContestSubmitAlready.Error(),
				}
				sseutil.SendMessage(ctx, &sseutil.ErrCodeFail, data)
				return
			}

			if err != nil {
				// Nếu lỗi là DB_NOT_FOUND, thêm người dùng vào contest
				if errors.Is(err, db.DB_NOT_FOUND) {
					if _, err := h.UserContestRepo.AddUserContest(ctx, db.AddUserContestParams{
						UserID:    claims.UserID,
						ContestID: contestInt,
						Questions: []byte(contestInfo.Questions),
					}); err != nil {
						// Ghi lỗi vào log với thông tin chi tiết
						slog.ErrorContext(ctx, "Error adding user_contest", "user_id", claims.UserID, "contest_id", contestInt, "error", err)

						// Trả lỗi cho client
						response.Error(ctx, err)
						return
					}
				} else {
					// Nếu lỗi khác, trả về lỗi ngay lập tức
					slog.ErrorContext(ctx, "Error getting user_contest", "user_id", claims.UserID, "contest_id", contestInt, "error", err)
					response.Error(ctx, err)
					return
				}
			}
			// Nếu không có lỗi, người dùng đã tham gia contest, không cần làm gì thêm
		}

		// Tạo client và đăng ký với hub
		client := channel.NewClient(ctx, claims.UserID, userInfo.Nickname, bufferSize)
		registerClientAndSendMessages(hub, client, contestInfo)

		// Tạo một context mới với khả năng hủy
		reqCtx, cancel := context.WithCancel(ctx.Request.Context())
		reqCtx = context.WithValue(reqCtx, "cancel", cancel)

		select {
		case <-reqCtx.Done():
			handleDisconnect(hub, client)
		case <-hub.Quit():
			h.Rooms.Remove(contestInt)
			return
		}

	}
}
