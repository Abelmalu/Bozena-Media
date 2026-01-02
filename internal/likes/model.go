package likes


type Like struct{
    Id  int `json:"id" db:"id"`
	UserID int `json:"user_id" db:"user_id" validate:"required,gt=0"`
	PostID int `json:"post_id"  db:"post_id" validate:"required,gt=0"`

	
}