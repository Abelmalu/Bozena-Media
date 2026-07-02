package dto
 


type RegisterResponse struct {

	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	

}


type LoginResponse struct {

	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	UserName 	string  `json:"username"`
	FollowerCount int 	`json:"follower_count"`
	FollowingCount int 	`json:"following_count"`
	ID 				int `json:"id"`


}