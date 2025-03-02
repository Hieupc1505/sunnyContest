package db

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
)

// Custom MarshalJSON để chuyển đổi `tags` thành []string khi xuất JSON
func (s SfSubject) MarshalJSON() ([]byte, error) {
	// Chuyển đổi tags từ pgtype.Text thành []string nếu có giá trị
	var tags []string
	if s.Tags.Valid {
		if err := json.Unmarshal([]byte(s.Tags.String), &tags); err != nil {
			return nil, err
		}
	}

	// Chuyển đổi toàn bộ struct sang JSON
	type Alias SfSubject
	return json.Marshal(&struct {
		Tags []string `json:"tags"`
		*Alias
	}{
		Tags:  tags,
		Alias: (*Alias)(&s),
	})
}

// Custom UnmarshalJSON để chuyển đổi `tags` từ []string thành pgtype.Text khi nhận dữ liệu
func (s *SfSubject) UnmarshalJSON(data []byte) error {
	type Alias SfSubject
	aux := &struct {
		Tags []string `json:"tags"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Chuyển đổi từ []string sang pgtype.Text
	tagsJSON, err := json.Marshal(aux.Tags)
	if err != nil {
		return err
	}

	s.Tags = pgtype.Text{String: string(tagsJSON), Valid: true}
	return nil
}

type AnswerItem struct {
	IsCorrect bool   `json:"is_correct"`
	Ans       string `json:"ans"`
}

// MarshalJSON ghi đè phương thức MarshalJSON để xử lý trường Answers
func (q SfQuestion) MarshalJSON() ([]byte, error) {
	// Định nghĩa một struct tạm thời để marshaling
	type Alias SfQuestion

	// Parse Answers từ string sang []AnswerItem
	var answerItems []AnswerItem
	if err := json.Unmarshal([]byte(q.Answers), &answerItems); err != nil {
		return nil, fmt.Errorf("failed to unmarshal answers: %v", err)
	}

	// Trả về JSON với Answers là []AnswerItem
	return json.Marshal(&struct {
		Answers []AnswerItem `json:"answers"`
		Alias
	}{
		Answers: answerItems,
		Alias:   (Alias)(q),
	})
}
