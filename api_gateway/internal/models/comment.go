package models

type NewComment struct {
	Author      string `json:"author"`
	CommentText string `json:"comment_text"`
	NewsId      int    `json:"news_id"`
	ParentId    int    `json:"parent_id,omitempty"`
}
