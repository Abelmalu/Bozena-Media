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
	"github.com/abelmalu/golang-posts/Feed/internal/dto"
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
	Image   string `json:"image"`
}

type follow struct {
	ID          int `json:"id"`
	FollowerID  int `json:"follower_id"`
	FollowingID int `json:"following_id"`
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
func StartConsumer(brokers []string, userTopic string, postTopic string, likedTopic, unLikedTopic string,followedTopic string, feedService core.FeedService, logger *platform.Logger) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating master consumer: %v", err)
	}

	log.Printf("Consumer started. Listening on topic: %s and %s", userTopic, postTopic)

	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		defer wg.Done()
		consumeUserEvents(consumer, userTopic, feedService, logger)
	}()

	go func() {
		defer wg.Done()
		consumePostEvents(consumer, postTopic, feedService, logger)
	}()

	go func() {
		defer wg.Done()
		postLikedConsumer(consumer, likedTopic, feedService, logger)
	}()

	go func() {
		defer wg.Done()
		postUnlikedConsumer(consumer, unLikedTopic, feedService, logger)
	}()

	go func() {
		defer wg.Done()
		followedConsumer(consumer, followedTopic, feedService, logger)
	}()

	wg.Wait()
}

func consumeUserEvents(consumer sarama.Consumer, topic string, feedService core.FeedService, logger *platform.Logger) {
	pc, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error consuming user topic: %v", err)
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

	jobs := make(chan *sarama.ConsumerMessage, 100)

	wg := sync.WaitGroup{}

	for range 100 {

		wg.Add(1)
		go workers(jobs, &wg, feedService, logger)

	}

	for {
		select {
		case msg := <-pc.Messages():

			jobs <- msg

		case err := <-pc.Errors():
			log.Printf("Post consumer error: %v", err)

		}
	}
}

func workers(msgs <-chan *sarama.ConsumerMessage, wgs *sync.WaitGroup, feedService core.FeedService, logger *platform.Logger) {

	defer wgs.Done()

	followCleint := initFollowClient()

	for msg := range msgs {

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

		var post PostCreatedPayload
		if err := json.Unmarshal(msg.Value, &post); err != nil {
			log.Printf("Failed to unmarshal post: %v", err)
			cancel()
			continue
		}

		wg := sync.WaitGroup{}

		wg.Add(2)

		go func() {
			defer wg.Done()

			err := feedService.CreateCachePost(ctx, post.ID, post.Title, post.Content, post.Image)

			if err != nil {
				logger.Error("Error inserting to posts_cache", zap.Error(err))
			}

		}()

		md := metadata.Pairs(
			"request-id", "askdfjalksdjfalsdkjflskdjf",
		)
		ctx = metadata.NewOutgoingContext(ctx, md)

		go func() {

			defer wg.Done()

			resp, err := followCleint.GetUserFollowers(
				ctx,
				&pb.GetUserFollowersRequest{
					FollowingId: int64(post.UserID),
					Limit:       int64(10),
				},
			)

			if err != nil {

				logger.Error("Erroor while getting followers from follow service ", zap.Error(err))

				return
			}

			followers := make([]int, 0, len(resp.Followers))

			for _, follower := range resp.Followers {

				user := UserCreatedPayload{
					ID: int(follower.UserId),
				}

				followers = append(followers, user.ID)

			}

			err = feedService.CreateFeedEntries(ctx, followers, post.ID, post.UserID)

			if err != nil {

				logger.Error("Error", zap.Error(err))
			}

		}()

		wg.Wait()

		cancel()

	}

}

func postLikedConsumer(consumer sarama.Consumer, postLikedTopic string, feedService core.FeedService, logger *platform.Logger) {

	pc, err := consumer.ConsumePartition(postLikedTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error consuming post partition: %v", err)
	}
	defer pc.Close()

	for {

		select {

		case msg := <-pc.Messages():

			var likeEvent struct {
				ID     int `json:"id"`
				UserID int `json:"user_id"`
				PostID int `json:"post_id"`
			}

			err := json.Unmarshal(msg.Value, &likeEvent)
			if err != nil {
				log.Printf("Failed to unmarshal user data: %v", err)
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

			if err := feedService.IncreaseLikeCount(ctx, likeEvent.PostID); err != nil {

				logger.Error("Error increasing like count ", zap.Error(err))
			}

			cancel()

		case err := <-pc.Errors():

			logger.Error("post liked consumer error", zap.Error(err))
		}
	}
}

func postUnlikedConsumer(consumer sarama.Consumer, postUnLikedTopic string, feedService core.FeedService, logger *platform.Logger) {

	pc, err := consumer.ConsumePartition(postUnLikedTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error consuming post partition: %v", err)
	}
	defer pc.Close()

	for {

		select {

		case msg := <-pc.Messages():

			var likeEvent struct {
				ID     int `json:"id"`
				UserID int `json:"user_id"`
				PostID int `json:"post_id"`
			}

			err := json.Unmarshal(msg.Value, &likeEvent)
			if err != nil {
				log.Printf("Failed to unmarshal user data: %v", err)
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

			if err := feedService.DecreaseLikeCount(ctx, likeEvent.PostID); err != nil {

				logger.Error("Error decreasing like count ", zap.Error(err))
			}

			fmt.Println("unlike coutn increased")
			cancel()

		case err := <-pc.Errors():

			logger.Error("post unliked consumer error", zap.Error(err))
		}
	}
}

func followedConsumer(consumer sarama.Consumer, followedTopic string, feedService core.FeedService, logger *platform.Logger) {

	pc, err := consumer.ConsumePartition(followedTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error consuming post partition: %v", err)
	}
	defer pc.Close()

	for {

		select {

		case msg := <-pc.Messages():

			var followed follow

			if err := json.Unmarshal(msg.Value, &followed); err != nil {

				logger.Error("Error while unmarshaling followed event", zap.Error(err))

				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

			resp, err := feedService.GetCachePosts(ctx, followed.FollowingID)

			if err != nil {
				logger.Error("Error getting cached posts", zap.Error(err))

				cancel()
				continue

			}

			feedEntries := make([]*dto.FeedEntry,3)

			for _,cachePost := range resp.CachePosts {

				feedEntry := dto.FeedEntry{

					OwnerID: cachePost.UserID,
					PostID: cachePost.PostID,
					UserID: followed.FollowerID,

				}

				feedEntries = append(feedEntries,&feedEntry)
			}

			if err := feedService.AddFeedEntries(ctx,feedEntries); err != nil {


				logger.Error("Error adding feed entries",zap.Error(err))
				continue
			}


			cancel()

		case err := <-pc.Errors():
			logger.Error("Post consumer error: %v", zap.Error(err))
		}
	}

}

// func unfollowedConsumer(consumer sarama.Consumer, unfollowedTopic string, feedService core.FeedService, logger *platform.Logger) {

// 	pc, err := consumer.ConsumePartition(unfollowedTopic, 0, sarama.OffsetNewest)
// 	if err != nil {
// 		log.Fatalf("Error consuming post partition: %v", err)
// 	}
// 	defer pc.Close()

// 	for {

// 		select {

// 		case msg := <-pc.Messages():

// 			var followed follow

// 			if err := json.Unmarshal(msg.Value, &followed); err != nil {

// 				logger.Error("Error while unmarshaling followed event", zap.Error(err))
// 				continue
// 			}

// 			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

// 			err := feedService.DecreaseFollowCounts(ctx, followed.FollowerID, followed.FollowingID)

// 			if err != nil {

// 				logger.Error("Error when decreasing follow count", zap.Error(err))
// 			}

// 			cancel()

// 		case err := <-pc.Errors():
// 			logger.Error("Post consumer error: %v", zap.Error(err))
// 		}
// 	}

// }
