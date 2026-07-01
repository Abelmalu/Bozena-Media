package models


type Follow struct {
	
	ID int `json:"id"`
	FollowerID int `json:"follower_id"`
	FollowingID int `json:"following_id"`
}


type User struct {


	ID int `json:"id"`
	Name string `json:"name"`
	Username string `json:"username"`

}