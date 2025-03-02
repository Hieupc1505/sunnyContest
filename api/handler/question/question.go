package question

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/ansutil"
	"go-rest-api-boilerplate/api/handler"
	"go-rest-api-boilerplate/api/middleware"
	"go-rest-api-boilerplate/api/response"
	app "go-rest-api-boilerplate/internal"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/internal/services/question"
	"log/slog"
	"strconv"
	"strings"
)

// Routes
func Routes(h *handler.Handler) {
	questionGroup := h.App.Group("/api/v1/web")
	questionSGroup := questionGroup.Use(middleware.AuthMiddleware(h))
	{
		questionSGroup.POST("/question/create", Create(h))
		questionSGroup.PUT("/question", Update(h))
		questionSGroup.GET("/questions", GetList(h))
		questionSGroup.GET("/questions/total", GetTotal(h))
		questionSGroup.GET("/question/:id", GetByID(h))
		questionSGroup.PUT("/question/delete/:id", Delete(h))
	}
}

func Create(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body question.QuestionParams
		if err := ctx.ShouldBindJSON(&body); err != nil {
			response.Error(ctx, nil)
			return
		}

		if body.QuestionImage != "" {
			result, err := h.ImgUploader.Upload(body.QuestionImage)
			if err != nil {
				slog.Info("Error uploading Image: %v", err)
				response.Error(ctx, app.ErrInternalServerError)
				return
			}
			body.QuestionImage = result.Url
		}

		ImgHandler := ansutil.GetHandler(body.AnswerType)
		if ImgHandler == nil {
			response.Error(ctx, app.ErrInvalidData)
			return
		}

		ans, err := ImgHandler.Handle(body.Answers, h.ImgUploader)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		body.Answers = ans

		result, err := h.QuestionService.AddQuestion(ctx, body)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		response.Success(ctx, result)
	}
}

func Update(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Binding dữ liệu từ request
		var body question.QuestionParams
		if err := ctx.ShouldBindJSON(&body); err != nil {
			response.Error(ctx, nil)
			return
		}
		fmt.Println("bodyID::::", body.ID)

		// Upload ảnh câu hỏi (nếu có)
		if body.QuestionImage != "" && !strings.HasPrefix(body.QuestionImage, ansutil.ImgHost) {
			result, err := h.ImgUploader.Upload(body.QuestionImage)
			if err != nil {
				slog.Error("Error uploading question image:", err)
				response.Error(ctx, app.ErrInternalServerError)
				return
			}
			body.QuestionImage = result.Url
		}

		// Xử lý ảnh câu trả lời (nếu có)
		fmt.Println("Answer_Type:::", body.AnswerType)
		ImgHandler := ansutil.GetHandler(body.AnswerType)
		if ImgHandler == nil {
			response.Error(ctx, app.ErrInvalidData)
			return
		}

		ans, err := ImgHandler.Handle(body.Answers, h.ImgUploader)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		body.Answers = ans

		// Cập nhật câu hỏi trong cơ sở dữ liệu
		result, err := h.QuestionService.UpdateQuestion(ctx, body)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		// Trả về kết quả
		response.Success(ctx, result)
	}
}
func GetTotal(h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy và xử lý tham số từ request
		subjectIDStr := c.Query("sid")

		// Chuyển đ��i tham số và xử lý l��i
		subjectID, err := strconv.ParseInt(subjectIDStr, 10, 64)
		if err != nil {
			response.Error(c, app.ErrInvalidData)
			return
		}

		// Gọi service để lấy t��ng số câu h��i theo môn học
		total, err := h.QuestionRepo.GetTotalQuestion(c, subjectID)
		if err != nil {
			response.Error(c, app.ErrInternalServerError)
			return
		}

		// Trả về kết quả
		response.Success(c, gin.H{"total": total})
	}
}
func GetList(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Lấy và xử lý tham số từ request
		subjectIDStr := ctx.Query("sid")
		limitStr := ctx.Query("limit")
		pageStr := ctx.Query("page")

		// Chuyển đổi tham số và xử lý lỗi
		subjectID, err := strconv.ParseInt(subjectIDStr, 10, 64)
		if err != nil {
			response.Error(ctx, app.ErrInvalidData)
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			response.Error(ctx, app.ErrInvalidData)
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 0 {
			response.Error(ctx, app.ErrInvalidData)
			return
		}

		// Tính toán offset (với page bắt đầu từ 0)
		offset := page * limit

		// Tạo params để gọi repository
		params := db.GetQuestionBySubjectIDParams{
			SubjectID: subjectID,
			Limit:     int32(limit),
			Offset:    int32(offset),
		}

		// Gọi repository để lấy dữ liệu
		results, err := h.QuestionRepo.GetQuestionBySubjectID(ctx, params)
		if err != nil {
			response.Error(ctx, nil)
			return
		}

		// Trả về kết quả thành công
		response.Success(ctx, results)
	}
}
func GetByID(h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.Error(c, app.ErrInvalidData)
			return
		}
		result, err := h.QuestionRepo.GetQuestionByID(c, idInt)
		if err != nil {
			response.Error(c, nil)
			return
		}
		response.Success(c, result)

	}
}
func Delete(h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
