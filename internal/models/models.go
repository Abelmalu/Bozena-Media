package models

import "time"

// users model
type User struct{

	ID int `json:"id" db:"id" `
	Name string `json:"name" db:"name" validate:"required,min=2,max=30 "`
	Username string `json:"username" db:"username" validate:"required,min=2,max=30"`
	Email string `json:"email" db:"email" validate:"email,required"`
	Password string `json:"password" db:"password" binding:"required,min=8"`
	Role string `json:"role" db:"role"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

// posts model
type Post struct {
	Id      string `json:"id" db:"id"`
	Title   string `json:"title"  db:"title"validate:"min=3,max=30,required"`
	Content string `json:"content" db:"content" validate:"min=5"`
	UserID  int    `json:"user_id" db:"user_id" validate:"required"`
}


//likes model
type Like struct{
    Id  int `json:"id" db:"id"`
	UserID int `json:"user_id" db:"user_id" validate:"required,gt=0"`
	PostID int `json:"post_id"  db:"post_id" validate:"required,gt=0"`

	
}

// RefreshTokens model 
type RefreshToken struct {
    ID          int       `json:"id"`
    UserID      int       `json:"user_id"`
    TokenText   string    `json:"-"` // Never export the hash to JSON for security
    ClientType  string    `json:"client_type"` // 'web' or 'mobile'
    ExpiresAt   time.Time `json:"expires_at"`
    Revoked     bool      `json:"revoked"`
    CreatedAt   time.Time `json:"created_at"`
}



