package dto

import model "github.com/abelmalu/golang-posts/Auth/internal/models"



type PaginatedResponse struct {

	Users []*model.User
	Cursor string 
	HasNext bool 
}