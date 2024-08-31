package dto

type Post struct {
	ID      int64    `json:"id"`
	Content string   `json:"content"`
	Title   string   `json:"title"`
	UserID  int64    `json:"user_id"`
	Tags    []string `json:"tags"`
}
