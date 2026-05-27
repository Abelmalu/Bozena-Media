package initiator

import (
	"log"

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
	logger *platform.Logger
}



func NewClient(logger *platform.Logger) *Client{

	  	postConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {

		log.Fatalf("failed to connect to gRPC server: %v", err)

	}
	authConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {

		log.Fatalf("failed to connect to gRPC server: %v", err)

	}


	likeConn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {

		log.Fatalf("failed to connect to gRPC server: %v", err)

	}

	followConn, err := grpc.NewClient("localhost:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {

		log.Fatalf("failed to connect to gRPC server: %v", err)

	}

	

	


	return &Client{
		authClient:client.NewAuthClient(authConn) ,
		postClient: client.NewPostClient(postConn),
		likeClient: client.NewLikeClient(likeConn),
		followClient: client.NewFollowClient(followConn),
		logger: logger,

	}

}