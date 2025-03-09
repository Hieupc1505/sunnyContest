package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"go-rest-api-boilerplate/types"
	"time"
)

type UserInfo struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Role        int32     `json:"role"`
	Status      int32     `json:"status"`
	CreatedTime time.Time `json:"created_time"`
	Profile     AddProfileRow
}

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

// MarshalJSON ghi đè phương thức MarshalJSON để xử lý trường Answers
func (q SfQuestion) MarshalJSON() ([]byte, error) {
	// Định nghĩa một struct tạm thời để marshaling
	type Alias SfQuestion

	// Parse Answers từ string sang []AnswerItem
	var answerItems []types.AnswerItem
	if err := json.Unmarshal([]byte(q.Answers), &answerItems); err != nil {
		return nil, fmt.Errorf("failed to unmarshal answers: %v", err)
	}

	// Trả về JSON với Answers là []AnswerItem
	return json.Marshal(&struct {
		Answers []types.AnswerItem `json:"answers"`
		Alias
	}{
		Answers: answerItems,
		Alias:   (Alias)(q),
	})
}

// UnmarshalJSON Phương thức để parse Answers từ chuỗi JSON thành []AnswerItem
func (q *SfQuestion) UnmarshalJSON(data []byte) error {
	// Định nghĩa struct tạm thời
	type Alias SfQuestion
	aux := &struct {
		*Alias
		Answers json.RawMessage `json:"answers"` // Sử dụng json.RawMessage để giữ nguyên dữ liệu JSON gốc
	}{
		Alias: (*Alias)(q),
	}

	// Unmarshal vào struct tạm
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("failed to unmarshal SfQuestion: %v", err)
	}

	// Chuyển đổi json.RawMessage thành string và gán vào trường Answers
	q.Answers = string(aux.Answers)

	return nil
}

// MarshalJSON overrides the MarshalJSON method to process Answers field
func (q GetUserByIDRow) MarshalJSON() ([]byte, error) {
	// Define an alias for GetUserByIDRow to avoid recursion in MarshalJSON
	type Alias UserInfo

	// Create a new Alias struct and remove Nickname and Avatar fields before marshaling
	info := UserInfo{
		ID:          q.ID,
		Username:    q.Username,
		Role:        q.Role,
		Status:      q.Status,
		CreatedTime: q.CreatedTime,
		Profile: AddProfileRow{
			Nickname: q.Nickname,
			Avatar:   q.Avatar,
		},
	}

	// Return the marshaled JSON with Profile and the rest of the fields (except Nickname and Avatar)
	return json.Marshal(&struct {
		Alias // This will include all fields from GetUserByIDRow except Nickname and Avatar
	}{
		Alias: (Alias)(info),
	})
}

// MarshalJSON triển khai custom marshalling cho SfContest
func (c SfContest) MarshalJSON() ([]byte, error) {
	// Parse trường Questions từ chuỗi JSON thành đối tượng Go
	var parsedQuestions interface{}
	if c.Questions != "" {
		fmt.Println("Handler marshall sfContest")
		if err := json.Unmarshal([]byte(c.Questions), &parsedQuestions); err != nil {
			return nil, errors.New("failed to unmarshal Questions: " + err.Error())
		}
	}

	// Tạo một struct tạm thời để marshalling
	type Alias SfContest
	return json.Marshal(&struct {
		Alias
		ParsedQuestions interface{} `json:"questions"`
	}{
		Alias:           (Alias)(c),
		ParsedQuestions: parsedQuestions,
	})
}

type Response struct {
	Error int         `json:"e"`
	Data  interface{} `json:"d,omitempty"`
}
