package auth


type User struct{

	Id int `json:"id" db:"id" `
	Name string `json:"name" db:"name" validate:"required,min=2,max=30 "`
	Username string `json:"username" db:"username" validate:"required,min=2,max=30"`
	Email string `json:"email" db:"email" validate:"email,required"`
	Password string `json:"password" db:"password" validate:"required,min=8"`
}


