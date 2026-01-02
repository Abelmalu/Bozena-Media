package auth


type User struct{

	Id int `json:"id" db:"id`
	name string `json:"name" db:"name"`
	username string `json:"username db:"username"`
	password string `json:password db:"password"`
}


