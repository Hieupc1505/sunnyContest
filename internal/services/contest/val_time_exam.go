package contest

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type TimeExam int

func (t TimeExam) GetTime() pgtype.Timestamptz {
	// Dereference t to get the actual int value of TimeExam
	timeValue := time.Now().Add(time.Duration(t) * time.Minute)

	return pgtype.Timestamptz{
		Time:  timeValue, // This is a time.Time value
		Valid: true,
	}
}

func (t TimeExam) parseInt() int32 {
	return int32(t)
}

func newTimeExam(t int32) (TimeExam, error) {
	if t == 0 {
		return 0, fmt.Errorf("timeExam is not a valid")
	}
	return TimeExam(t), nil
}
