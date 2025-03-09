package user_contest

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/timeutil"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/types"
	"math"
	"time"
)

func MakeResult(userAnswers types.UserSubmitBody, questions []types.Question, timeStart time.Time) (types.Results, []types.QuestionItemInfo) {
	correctCount := 0
	var result types.Results
	var infos []types.QuestionItemInfo

	result.Results = make([]types.FormattedResult, 0)
	infos = make([]types.QuestionItemInfo, 0)

	// Tạo một map để lưu trữ câu hỏi theo ID để truy cập nhanh
	answersMap := make(map[int64]types.SubnitItem)
	for _, ans := range userAnswers.Answers {
		answersMap[ans.QuestionID] = ans
	}

	// Duyệt qua các câu trả lời của người dùng
	for _, q := range questions {
		var correct types.ExamInfo
		var info types.QuestionItemInfo
		exam := &types.ExamInfo{
			Index:     IndexUserAnswerMiss,
			Ans:       "",
			IsCorrect: false,
		} // Khai báo exam là nil ban đầu

		// Lấy câu hỏi tương ứng từ map
		userAnswer, exists := answersMap[q.ID]

		info.Exam = IndexUserAnswerMiss // Gán giá trị mặc định trước khi vòng lặp bắt đầu

		for key, ans := range q.Answers {
			info.ID = userAnswer.QuestionID
			if exists && key == userAnswer.Index {
				// Khởi tạo exam chỉ khi điều kiện đúng
				exam.Index = key
				exam.Ans = ans.Ans
				exam.IsCorrect = ans.IsCorrect
				info.Exam = key
			}
			if ans.IsCorrect {
				correct.Index = key
				correct.Ans = ans.Ans
				correct.IsCorrect = ans.IsCorrect
				info.Correct = key
			}
		}
		if exam != nil && exam.IsCorrect {
			correctCount++
		}
		resultItem := types.FormattedResult{
			QuestionID: q.ID,
			Exam:       exam,
			Correct:    correct,
		}
		// Thêm vào kết quả
		result.Results = append(result.Results, resultItem)
		infos = append(infos, info)
	}
	result.NumCorrect = correctCount
	result.NumIncorrect = len(questions) - correctCount
	result.TimeSubmit = math.Round(timeutil.SinceFromToSecond(timeStart)*1000) / 1000

	return result, infos
}

// SubmitContest to handler when user submit
func (s *Service) SubmitContest(ctx *gin.Context, contestID, userID int64, timeStart time.Time, userAnswers types.UserSubmitBody, questions []types.Question) (*types.Results, error) {

	results, infos := MakeResult(userAnswers, questions, timeStart)

	examJson, err := json.Marshal(infos)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal exam: %v", err)
	}

	resultsInfo := NewResults(results.NumCorrect, results.NumIncorrect, results.TimeSubmit)
	resultsInfoJson, err := json.Marshal(resultsInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal rank: %v", err)
	}

	if err := s.repo.UpdateExamAndResult(ctx, db.UpdateExamAndResultParams{
		Exam:      examJson,        // examJson là json.RawMessage hoặc nil
		Result:    resultsInfoJson, // rankJson là json.RawMessage hoặc nil
		ContestID: contestID,
		UserID:    userID,
	}); err != nil {
		return nil, fmt.Errorf("failed to update exam and result: %v", err)
	}

	return &results, nil

}
