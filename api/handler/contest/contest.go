package contest

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/contextutil"
	"go-rest-api-boilerplate/api/handler"
	"go-rest-api-boilerplate/api/middleware"
	"go-rest-api-boilerplate/api/response"
	"go-rest-api-boilerplate/api/sseutil"
	app "go-rest-api-boilerplate/internal"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/internal/services/contest"
	"go-rest-api-boilerplate/internal/services/question"
	"go-rest-api-boilerplate/pkg/sse/channel"
	"go-rest-api-boilerplate/types"
	"log/slog"
	"strconv"
)

func Routes(h *handler.Handler) {
	ContestGroup := h.App.Group("/api/v1/contest")
	ContestSGroup := ContestGroup.Use(middleware.AuthMiddleware(h))
	{
		ContestSGroup.POST("", AddAndUpdate(h))
		ContestSGroup.GET("/start/:contestId", StartContest(h))
		ContestSGroup.GET("/stop/:contestId", EndContest(h))
		ContestSGroup.GET("/live/me", MyLiveContest(h))
		ContestSGroup.GET("/live", GetLiveContest(h))
		ContestSGroup.POST("/:contestId/submit-paper", SubmitContest(h))
		ContestSGroup.GET("/live/delete/:id", nil)
		ContestSGroup.GET("/live/:id", GetDetailLiveContest(h))
		ContestSGroup.POST("/play/:contestId", LiveContest(h))
		ContestSGroup.GET("/me", TeacherContest(h))
		ContestSGroup.GET("/joins", JoinedContest(h))
		ContestSGroup.GET("/:id", Statistics(h))
	}
}

func JoinedContest(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := contextutil.GetUser(ctx)
		if err != nil {
			response.Error(ctx, app.ErrInternalServerError)
		}
		data, err := h.UserContestRepo.GetUserContestsJoined(ctx, userId)
		if err != nil {
			response.Error(ctx, app.ErrInternalServerError)
			return
		}
		response.Success(ctx, data)
	}
}

func TeacherContest(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := contextutil.GetUser(ctx)
		if err != nil {
			response.Error(ctx, app.ErrInternalServerError)
			return
		}

		limitStr := ctx.Query("limit")
		pageStr := ctx.Query("page")

		// Đặt giá trị mặc định cho limit và page
		limit := 10
		page := 0

		// Nếu có giá trị limit, chuyển đổi từ string sang int
		if limitStr != "" {
			parsedLimit, err := strconv.Atoi(limitStr)
			if err != nil {
				response.Error(ctx, app.ErrInvalidData)
				return
			}
			limit = parsedLimit
		}

		// Nếu có giá trị page, chuyển đổi từ string sang int
		if pageStr != "" {
			parsedPage, err := strconv.Atoi(pageStr)
			if err != nil || parsedPage < 0 {
				response.Error(ctx, app.ErrInvalidData)
				return
			}
			page = parsedPage
		}

		// Tính toán offset (với page bắt đầu từ 0)
		offset := page * limit

		// Tạo params để gọi repository
		params := db.GetContestsForTeacherParams{
			UserID: userId,
			Limit:  int32(limit),
			Offset: int32(offset),
		}

		contests, err := h.ContestRepo.GetContestsForTeacher(ctx, params)
		if err != nil {
			slog.Error("Get teacher contests error", "error", err)
			response.Error(ctx, app.ErrInternalServerError)
			return
		}

		response.Success(ctx, contests)
	}
}

func Statistics(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		contestId := ctx.Param("id")
		contestInt, err := strconv.ParseInt(contestId, 10, 64)
		if err != nil {
			response.Error(ctx, app.ErrInternalServerError)
			return
		}

		contestInfo, err := h.ContestRepo.GetContestDetailByID(ctx, contestInt)
		if err != nil {
			if errors.Is(err, db.DB_NOT_FOUND) {
				response.Error(ctx, app.ErrContestNotFound)
				return
			}
			slog.Error("Get contest statistics error", "error", err)
		}

		users, err := h.ContestRepo.GetUsersInContest(ctx, contestInt)
		if err != nil {
			slog.Error("Get contest users error", "error", err)
			response.Error(ctx, app.ErrInternalServerError)
		}

		var questions []db.SfQuestion
		err = json.Unmarshal([]byte(contestInfo.Questions), &questions)
		if err != nil {
			slog.Error("Marshal questions error", "error", err)
			response.Error(ctx, app.ErrInternalServerError)
			return
		}

		data := contest.MakeStatisticsRsp(contestInfo, users, questions)

		response.Success(ctx, data)

	}
}

func GetDetailLiveContest(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		contestId := ctx.Param("id")
		contestInt, err := strconv.ParseInt(contestId, 10, 64)
		if err != nil {
			response.Error(ctx, app.ErrInternalServerError)
			return
		}

		result, err := h.ContestRepo.GetContestByID(ctx, contestInt)
		if err != nil {
			if errors.Is(err, db.DB_NOT_FOUND) {
				response.Error(ctx, app.ErrContestNotFound)
				return
			}
			slog.Error("Get contest live error", "error", err)
		}

		response.Success(ctx, result)
	}
}

func GetLiveContest(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		contests, err := h.ContestRepo.GetListLiveContest(ctx)
		if err != nil {
			slog.Error("Get list live contest error", "error", err)
			response.Error(ctx, app.ErrInternalServerError)
			return
		}
		response.Success(ctx, contests)
	}
}

func SubmitContest(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := contextutil.GetUser(ctx)
		if err != nil {
			response.Error(ctx, app.ErrInternalServerError)
		}

		contestId := ctx.Param("contestId")
		contestInt, err := strconv.ParseInt(contestId, 10, 64)
		if err != nil {
			response.Error(ctx, app.ErrInvalidData)
			return
		}

		var submit types.UserSubmitBody
		if err := ctx.ShouldBindJSON(&submit); err != nil {
			slog.Error("Invalid parse data submit contest", "error", err)
			response.Error(ctx, app.ErrInvalidData)
			return
		}

		contestInfo, err := h.ContestRepo.GetContestByID(ctx, contestInt)
		if err != nil {
			slog.Error("Get contest live error", "error", err)
			response.Error(ctx, app.ErrInternalServerError)
			return
		}

		// Lấy hub tương ứng với contest
		hub, err := h.Rooms.GetHub(contestInt)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		questions := hub.GetQuestions()

		results, err := h.UserContestService.SubmitContest(ctx, contestInt, userId, contestInfo.TimeStartExam.Time, submit, questions)
		if err != nil {
			slog.Error("Submit contest error::::", "error", err)
			response.Error(ctx, app.ErrInternalServerError)
			return
		}

		if err := hub.UserSubmit(userId, results); err != nil {
			slog.Error("User submit error", "error", err)
			response.Error(ctx, app.ErrInternalServerError)
			return
		}

		hub.Broadcast(sseutil.NewSseRes(types.UserJoin, hub.GetUsers(contestInfo.UserID)))

		response.Success(ctx, gin.H{"success": true})

		// Kiểm tra xem người dùng đã đăng nhập vào contest chưa
	}
}

func EndContest(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		contestId := ctx.Param("contestId")
		contestInt, err := strconv.ParseInt(contestId, 10, 64)
		if err != nil {
			response.Error(ctx, app.ErrInvalidData)
			return
		}

		// Lấy hub tương ứng với contest
		hub, err := h.Rooms.GetHub(contestInt)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		// Xóa contest từ room
		//hub.EndContest()
		//hub.Stop("Route end contest")

		hub.StopTimer()

		response.Success(ctx, gin.H{"id": contestInt})
	}
}

func StartContest(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		contestId := ctx.Param("contestId")
		contestInt, err := strconv.ParseInt(contestId, 10, 64)
		if err != nil {
			response.Error(ctx, app.ErrInvalidData)
			return
		}

		// Lấy hub tương ứng với contest
		hub, err := h.Rooms.GetHub(contestInt)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		contestInfo, err := h.ContestRepo.GetContestLiveByID(ctx, contestInt)
		if err != nil {
			slog.Error("Get contest live error", "error", err)
			response.Error(ctx, app.ErrInternalServerError)
			return
		}

		// Kích hoạt bộ đếm thời gian
		go hub.StartTimer(contestInfo.TimeExam)

		if err := h.ContestRepo.StartContest(ctx, contestInt); err != nil {
			response.Error(ctx, app.ErrInternalServerError)
			return
		}

		response.Success(ctx, gin.H{"id": contestInt})
	}
}

// LiveContest live contest for member can access
func LiveContest(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		contestId := ctx.Param("contestId")
		contestInt, err := strconv.ParseInt(contestId, 10, 64)
		if err != nil {
			response.Error(ctx, app.ErrInvalidData)
			return
		}

		contestInfo, err := h.ContestRepo.GetContestLiveByID(ctx, contestInt)
		if err != nil {
			slog.Error("Get contest live error", "error", err)
			response.Error(ctx, app.ErrInternalServerError)
			return
		}

		params := question.NewRandomQuestionParams(contestInfo.SubjectID, contestInfo.NumQuestion)
		questions, err := h.QuestionRepo.GetRandomQuestions(ctx, *params)
		if err != nil {
			slog.Error("Get random questions error", "error", err)
			response.Error(ctx, app.ErrInternalServerError)
			return
		}

		questionJson, err := json.Marshal(questions)
		if err != nil {
			slog.Error("Marshal questions error", "error", err)
			response.Error(ctx, app.ErrInternalServerError)
			return
		}

		if err := h.ContestRepo.UpdateContestStateAndQuestions(ctx, db.UpdateContestStateAndQuestionsParams{
			ID:        contestInt,
			State:     db.ContestStateWAITING,
			Questions: string(questionJson),
		}); err != nil {
			slog.Error("Update contest state", "error", err)
			response.Error(ctx, app.ErrInternalServerError)
			return
		}

		hub := channel.NewHub(h, contestInt)
		hub.Publish()
		hub.SetQuestions(questions)
		h.Rooms.Add(contestInt, hub)

		response.Success(ctx, gin.H{"id": contestInt})

	}
}

func MyLiveContest(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := contextutil.GetUser(ctx)
		if err != nil {
			response.Error(ctx, app.ErrInternalServerError)
			return
		}
		context, err := h.ContestRepo.GetMyContestLive(ctx, userId)
		if err != nil {
			if errors.Is(err, db.DB_NOT_FOUND) {
				response.Success(ctx, nil)
				return
			}
			slog.Error("Get my contest live error", "error", err)
			response.Error(ctx, app.ErrInternalServerError)
			return
		}
		response.Success(ctx, context)
	}
}

func AddAndUpdate(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data contest.AddAndUpdateParams
		if err := ctx.ShouldBindJSON(&data); err != nil {
			response.Error(ctx, app.ErrInvalidData)
			return
		}

		var cont *db.SfContest
		var err error

		totalQ, err := h.QuestionRepo.GetTotalQuestion(ctx, data.SubjectID)
		if err != nil {
			slog.Error("Get total question error", "error", err)
			response.Error(ctx, err)
			return
		}
		if data.ID == 0 {
			cont, err = h.ContestService.Add(ctx, data, totalQ)
			if err != nil {
				slog.Error("Create contest error", "error", err)
				response.Error(ctx, err)
				return
			}
		} else {
			cont, err = h.ContestService.Update(ctx, data, totalQ)
			if err != nil {
				slog.Error("Update contest error", "error", err)
				response.Error(ctx, err)
				return
			}
		}
		response.Success(ctx, gin.H{"id": cont.ID})
	}
}
