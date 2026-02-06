package model

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


// client types for mobile apps and browsers
type ClientType string 

const (
	ClientWeb ClientType = "web"
	ClientMobile ClientType = "mobile"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
