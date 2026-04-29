package initiator

import (
	"github.com/abelmalu/golang-posts/APIGateway/internal/handler"
)


type Handler struct{

	authHandler handler.AuthHandler 
	postHandler handler.PostHandler


}

func InitHandler(client Client ) *Handler{


	return &Handler{
		authHandler: *handler.NewAuthHandler(&client.authClient),
		postHandler: *handler.NewPostHandler(&client.postClient),
	}


	


}