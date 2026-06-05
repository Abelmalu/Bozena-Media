package models


type Follow struct {
	
	ID int 
	FollowerID int 
	FollowingID int
}


type User struct {


	ID int `json:"id"`
	Name string `json:"name"`
	Username string `json:"username"`

}