package contest

import db "go-rest-api-boilerplate/internal/db/sqlc"

type AddAndUpdateParams struct {
	ID          int64  `json:"id"`
	SubjectID   int64  `json:"subject_id"`
	NumQuestion int32  `json:"num_question"`
	TimeExam    int32  `json:"time_exam"`
	SubjectName string `json:"subject_name"`
}

func (p AddAndUpdateParams) ToAddParam() db.CreateContestParams {
	return db.CreateContestParams{
		UserID:      p.ID,
		SubjectID:   p.SubjectID,
		NumQuestion: p.NumQuestion,
		TimeExam:    p.TimeExam,
		State:       db.ContestStateRUNNING,
	}
}

func NewAddContest(userID int64, subjectID int64, numQuestion NumQuestion, timeExam TimeExam) *db.CreateContestParams {
	return &db.CreateContestParams{
		UserID:      userID,
		SubjectID:   subjectID,
		NumQuestion: numQuestion.Int(),
		TimeExam:    timeExam.parseInt(),
		State:       db.ContestStateIDLE,
	}
}

func NewUpdateContest(id, userId, subjectID int64, numQuestion NumQuestion, timeExam TimeExam) *db.UpdateContestParams {
	return &db.UpdateContestParams{
		ID:          id,
		UserID:      userId,
		SubjectID:   subjectID,
		NumQuestion: numQuestion.Int(),
		TimeExam:    timeExam.parseInt(),
		State:       db.ContestStateIDLE,
	}
}

// ToSfContest convert *db.GetContestLiveByIDRow to db.SfContest
func ToSfContest(s *db.GetContestLiveByIDRow) *db.SfContest {
	return &db.SfContest{
		ID:          s.ID,
		UserID:      s.UserID,
		SubjectID:   s.SubjectID,
		NumQuestion: s.NumQuestion,
		TimeExam:    s.TimeExam,
		Questions:   s.Questions,
		State:       s.State,
	}
}
