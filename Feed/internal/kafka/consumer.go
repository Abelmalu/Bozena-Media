package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/abelmalu/golang-posts/Feed/internal/core"
	"github.com/abelmalu/golang-posts/follow/proto/pb"
	"github.com/abelmalu/golang-posts/platform"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Define structures for your events
type UserCreatedPayload struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type PostCreatedPayload struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int    `json:"user_id"`
}

func initFollowClient() pb.FollowServiceClient {

	followConn, err := grpc.NewClient("localhost:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {

		log.Fatalf("failed to connect to gRPC server: %v", err)

	}

	followClient := pb.NewFollowServiceClient(followConn)

	return followClient
}

// StartEventConsumers initializes separate listeners for different topics
func StartConsumer(brokers []string, userTopic string, postTopic string, feedService core.FeedService, logger *platform.Logger) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating master consumer: %v", err)
	}

	log.Printf("Consumer started. Listening on topic: %s and %s", userTopic, postTopic)

	var wg sync.WaitGroup
	wg.Add(2)

	// Launch User Consumer
	go func() {
		defer wg.Done()
		consumeUserEvents(consumer, userTopic, feedService, logger)
	}()

	// Launch Post Consumer
	go func() {
		defer wg.Done()
		consumePostEvents(consumer, postTopic, feedService, logger)
	}()

	wg.Wait()
}

func consumeUserEvents(consumer sarama.Consumer, topic string, feedService core.FeedService, logger *platform.Logger) {
	pc, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error consuming user partition: %v", err)
	}
	defer pc.Close()

	for {
		select {
		case msg := <-pc.Messages():
			var user UserCreatedPayload
			if err := json.Unmarshal(msg.Value, &user); err != nil {
				log.Printf("Failed to unmarshal user: %v", err)
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			err = feedService.CreateCacheUser(ctx, user.ID, user.Username, user.Name)
			cancel() // call immediately instead of defer in a loop to avoid memory leaks

			if err != nil {
				logger.Error("Error inserting to users_cache", zap.Error(err))
			}

		case err := <-pc.Errors():
			log.Printf("User consumer error: %v", err)
		}
	}
}

func consumePostEvents(consumer sarama.Consumer, topic string, feedService core.FeedService, logger *platform.Logger) {
	pc, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error consuming post partition: %v", err)
	}
	defer pc.Close()

	followCleint := initFollowClient()

	for {
		select {
		case msg := <-pc.Messages():
			var post PostCreatedPayload
			if err := json.Unmarshal(msg.Value, &post); err != nil {
				log.Printf("Failed to unmarshal post: %v", err)
				continue
			}

			ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
			// Assuming you implement CreateCachePost in your feedService
			err = feedService.CreateCachePost(ctx, post.ID, post.Title, post.Content)

			if err != nil {
				logger.Error("Error inserting to posts_cache", zap.Error(err))
			}

			md := metadata.Pairs(
				"request-id", "askdfjalksdjfalsdkjflskdjf",
			)
			ctx = metadata.NewOutgoingContext(ctx, md)

			resp, err := followCleint.GetUserFollowers(
				ctx,
				&pb.GetUserFollowersRequest{
					FollowingId: int64(post.UserID),
					Limit:       int64(10),
				},
			)

			if err != nil {

				logger.Error("Erroor while getting followers from follow service ", zap.Error(err))
			}

			followers := make([]int, 0, len(resp.Followers))

			for _, follower := range resp.Followers {

				user := UserCreatedPayload{
					ID: int(follower.UserId),
				}

				followers = append(followers, user.ID)

			}

			fmt.Println("*****************followers**************", followers)
		case err := <-pc.Errors():
			log.Printf("Post consumer error: %v", err)
		}
	}
}
