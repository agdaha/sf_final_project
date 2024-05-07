package models

import (
	"database/sql"
	"encoding/json"
)

type NewDetailed struct {
	News     News
	Comments Comments
}

type News struct {
	Id          int
	Title       string
	Description string
	Link        string
	PubDate     int
	Author      string
	Guid        string
}

type Comment struct {
	Id          int       `json:"id,omitempty"`
	Author      string    `json:"author"`
	CommentText string    `json:"comment_text"`
	NewsId      int       `json:"news_id"`
	ParentId    NullInt64 `json:"parent_id,omitempty"`
	Level       int       `json:"level,omitempty"`
}

type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON for NullInt64
func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// Unmarshalling into a pointer will let us detect null
func (v *NullInt64) UnmarshalJSON(data []byte) error {
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

type RoutineComments struct {
	Comments Comments
	Err      error
}

type RoutineNews struct {
	News News
	Err  error
}

type Comments []Comment
