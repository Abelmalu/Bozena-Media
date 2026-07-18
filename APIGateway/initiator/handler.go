package initiator

import (
	"github.com/abelmalu/golang-posts/APIGateway/internal/handler"
	"github.com/abelmalu/golang-posts/platform"
)

type Handler struct {
	ah                  *handler.AuthHandler
	ps                  *handler.PostHandler
	lh                  *handler.LikeHandler
	fh                  *handler.FollowHandler
	fd                  *handler.FeedHandler
	notificationHandler *handler.NotificationHanlder
}

func InitHandler(client Client, logger *platform.Logger) *Handler {

	return &Handler{
		ah:                  handler.NewAuthHandler(client.authClient, logger),
		ps:                  handler.NewPostHandler(client.postClient, logger),
		lh:                  handler.NewLikeHandler(client.likeClient, logger),
		fh:                  handler.NewFollowHandler(client.followClient, logger),
		fd:                  handler.NewFeedHandler(client.feedClient, logger),
		notificationHandler: handler.NewNotificationHandler(logger),
	}

}
