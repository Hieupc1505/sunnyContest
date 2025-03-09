package contest

import (
	"encoding/json"
	"fmt"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/internal/services/question"
	"go-rest-api-boilerplate/internal/services/user_contest"
	"go-rest-api-boilerplate/types"
	"log"
)

type QuestionsMap map[int64][]types.AnswerItem

func MakeResults(q QuestionsMap, userExam []types.QuestionItemInfo) []types.FormattedResult {
	var results []types.FormattedResult
	var item types.FormattedResult

	exam := &types.ExamInfo{}
	correct := &types.ExamInfo{}

	results = make([]types.FormattedResult, 0)
	for _, ans := range userExam {
		exam.Index = ans.Exam
		exam.IsCorrect = ans.Exam == ans.Correct
		if ans.Exam != user_contest.IndexUserAnswerMiss {
			exam.Ans = q[ans.ID][ans.Exam].Ans
		}
		correct.Index = ans.Correct
		correct.Ans = q[ans.ID][ans.Correct].Ans
		correct.IsCorrect = true

		item.QuestionID = ans.ID
		item.Exam = exam
		item.Correct = *correct

		results = append(results, item)
	}

	return results
}

func MakeUserItem(user db.GetUsersInContestRow, q QuestionsMap) (userItem *types.UserResultContest) {

	var userExam []types.QuestionItemInfo
	var results types.Results

	if err := json.Unmarshal(user.Exam, &userExam); user.Exam != nil && err != nil {
		fmt.Println(user.Exam)
		log.Fatal("Error unmarshalling user exam::::", err)
	}
	if err := json.Unmarshal(user.Result, &results); user.Result != nil && err != nil {
		fmt.Println(user.Result)
		log.Fatal("Error unmarshalling user result::::", err)

	}

	data := MakeResults(q, userExam)
	results.Results = data

	userItem = &types.UserResultContest{
		UserID:   user.ID,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Results:  results,
	}
	return
}

func MakeUsers(users []db.GetUsersInContestRow, q QuestionsMap) (data []types.UserResultContest) {
	for _, user := range users {
		rsp := MakeUserItem(user, q)
		data = append(data, *rsp)
	}
	return
}

// MakeStatisticsRsp returns statistics of contest for member and teacher
func MakeStatisticsRsp(c db.GetContestDetailByIDRow, users []db.GetUsersInContestRow, questions []db.SfQuestion) (data *types.StatisticResult) {

	q := make(QuestionsMap)
	for _, qs := range questions {
		q[qs.ID] = question.ConvertAnswersToJson(qs.Answers)
	}

	usersMap := MakeUsers(users, q)
	return &types.StatisticResult{
		ID:          c.ID,
		TimeExam:    c.TimeExam,
		NumQuestion: c.NumQuestion,
		State:       string(c.State),
		Questions:   questions,
		SubjectName: c.SubjectName,
		Users:       usersMap,
	}
}
