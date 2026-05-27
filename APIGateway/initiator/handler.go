package initiator

import (
	"github.com/abelmalu/golang-posts/APIGateway/internal/handler"
	"github.com/abelmalu/golang-posts/platform"
)


type Handler struct{

	authHandler *handler.AuthHandler 
	postHandler *handler.PostHandler
	likeHandler *handler.LikeHandler
	followHandler *handler.FollowHandler

}

func InitHandler(client Client,logger *platform.Logger ) *Handler{


	return &Handler{
		authHandler: handler.NewAuthHandler(client.authClient,logger),
		postHandler: handler.NewPostHandler(client.postClient,logger),
		likeHandler: handler.NewLikeHandler(client.likeClient,logger),
		followHandler: handler.NewFollowHandler(client.followClient,logger),
	}


	


}