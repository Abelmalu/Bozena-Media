package dto

import "time"

type UserNotification struct {
	ID   int    `json:"id" db:"id"`
	UseraName string    `json:"username" db:"username"`
	ActorID   int       `json:"actor_id" db:"actor_id"`
	Message   string    `json:"message" db:"message"`
	CreatedAT *time.Time `json:"created_at" db:"created_at"`
}

type PaginatedResponse struct {
	UserNotifications []*UserNotification
	Cursor            string
	HasNext           bool
}


type User struct {

	ID int
	UserName string 
	Name string 

}
