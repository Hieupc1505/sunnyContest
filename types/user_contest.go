package types

type AnswerItem struct {
	IsCorrect bool   `json:"is_correct"`
	Ans       string `json:"ans"`
}

type ExamAnswer struct {
	Index int `json:"index"`
}

type Results struct {
	NumCorrect   int               `json:"num_correct"`
	NumIncorrect int               `json:"num_incorrect"`
	TimeSubmit   float64           `json:"time_submit"`
	Results      []FormattedResult `json:"results,omitempty"`
}

type FormattedResult struct {
	QuestionID int64     `json:"question_id"`
	Exam       *ExamInfo `json:"exam"`
	Correct    ExamInfo  `json:"correct"`
}

type ExamInfo struct {
	Index     int    `json:"index,omitempty"`
	Ans       string `json:"ans,omitempty"`
	IsCorrect bool   `json:"is_correct,omitempty"`
}

type ContestItemResult struct {
	ID       int64    `json:"id"`
	Owner    bool     `json:"owner"`
	Nickname string   `json:"nickname"`
	Results  *Results `json:"results,omitempty"`
}

type SubnitItem struct {
	QuestionID int64 `json:"question_id"`
	Index      int   `json:"index"`
}

type UserSubmitBody struct {
	Answers []SubnitItem `json:"answers"`
}

type QuestionItemInfo struct {
	ID      int64 `json:"question_id"`
	Exam    int   `json:"exam_idx"`
	Correct int   `json:"correct_idx"`
}

type UserResultContest struct {
	UserID   int64   `json:"user_id"`
	Nickname string  `json:"nickname"`
	Avatar   string  `json:"avatar"`
	Results  Results `json:"results"`
}

type StatisticResult struct {
	ID          int64               `json:"id"`
	TimeExam    int32               `json:"time_exam"`
	NumQuestion int32               `json:"num_question"`
	State       string              `json:"state"`
	Questions   any                 `json:"questions"`
	SubjectName string              `json:"subject_name"`
	Users       []UserResultContest `json:"users,omitempty"`
}
