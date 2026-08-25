package initiator

import (
	"log"

	"github.com/abelmalu/golang-posts/APIGateway/config"
	client "github.com/abelmalu/golang-posts/APIGateway/internal/clients"
	"github.com/abelmalu/golang-posts/platform"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)



type Client struct {

	authClient *client.AuthClient
	postClient *client.PostClient
	likeClient *client.LikeClient
	followClient *client.FollowClient
	feedClient *client.FeedClient
	logger *platform.Logger
}



func NewClient(logger *platform.Logger,cfg *config.Config) *Client{

	  	postConn, err := grpc.NewClient(cfg.PostServiceADD, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {

		log.Fatalf("failed to connect to post gRPC server: %v", err)

	}
	authConn, err := grpc.NewClient(cfg.AuthServiceADD, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {

		log.Fatalf("failed to connect to auth gRPC server: %v", err)

	}


	likeConn, err := grpc.NewClient(cfg.LikeServiceADD, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {

		log.Fatalf("failed to connect to like gRPC server: %v", err)

	}

	followConn, err := grpc.NewClient(cfg.FollowServiceADD, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {

		log.Fatalf("failed to connect to follow gRPC server: %v", err)

	}
    

	feedConn, err := grpc.NewClient(cfg.FeedServiceADD, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {

		log.Fatalf("failed to connect to gRPC server: %v", err)

	}

	

	


	return &Client{
		authClient:client.NewAuthClient(authConn) ,
		postClient: client.NewPostClient(postConn),
		likeClient: client.NewLikeClient(likeConn),
		followClient: client.NewFollowClient(followConn),
		feedClient: client.NewFeedClient(feedConn),
		logger: logger,

	}

}