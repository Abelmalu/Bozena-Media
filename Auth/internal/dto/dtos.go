package dto

import model "github.com/abelmalu/golang-posts/Auth/internal/models"



type PaginatedResponse struct {

	Users []*model.User
	Cursor string 
	HasNext bool 
}

type UserProfileResponse struct {

	ID int64 `json:"id" db:"id"`
	UserName string `json:"user_name" db:"username"`
	Name 	 string `json:"name" db:"name"`
	Avatar *string `json:"avatar" db:"avatar"`

}