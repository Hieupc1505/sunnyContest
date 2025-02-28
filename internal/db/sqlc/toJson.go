package db

import (
	"encoding/json"
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
