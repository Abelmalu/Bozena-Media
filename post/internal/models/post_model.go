package models


// posts model
type Post struct {
	Id      string `json:"id" db:"id"`
	Title   string `json:"title"  db:"title" validate:"min=3,max=30,required"`
	Content string `json:"content" db:"content" validate:"min=5"`
	UserID  int    `json:"user_id" db:"user_id" validate:"required"`
}

