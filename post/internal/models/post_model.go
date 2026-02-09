package models


// posts model
type Post struct {
	ID      string `json:"id" db:"id"`
	Title   string `json:"title"  db:"title" validate:"min=3,max=30,required"`
	Content string `json:"content" db:"content" validate:"min=5"`
	UserID  int    `json:"user_id" db:"user_id" validate:"required"`
}

//likes model
type Like struct{
    Id  int `json:"id" db:"id"`
	UserID int `json:"user_id" db:"user_id" validate:"required,gt=0"`
	PostID int `json:"post_id"  db:"post_id" validate:"required,gt=0"`

	
}
