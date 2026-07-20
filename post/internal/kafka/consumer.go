package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/abelmalu/golang-posts/post/internal/core"
	"go.uber.org/zap"
)

func StartConsumer(brokers []string, userCreatedtopic, postLikedTopic, postUnlikedTopic string, postService core.PostService, logger *platform.Logger) {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating master consumer: %v", err)
	}

	go func() {

		go userCreatedConsumer(consumer, userCreatedtopic, postService, logger)
	}()
	go func() {

		go postLikedConsumer(consumer, postLikedTopic, postService, logger)
	}()

	go func() {

		go postUnlikedConsumer(consumer, postUnlikedTopic, postService, logger)
	}()

}

func userCreatedConsumer(consumer sarama.Consumer, userCreatedTopic string, postService core.PostService, logger *platform.Logger) {

	pc, err := consumer.ConsumePartition(userCreatedTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error consuming post partition: %v", err)
	}
	defer pc.Close()

	for {

		select {
		case msg := <-pc.Messages():
			var user struct {
				ID       int    `json:"id"`
				Username string `json:"username"`
				Name     string `json:"name"`
			}

			err := json.Unmarshal(msg.Value, &user)
			if err != nil {
				log.Printf("Failed to unmarshal user data: %v", err)
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

			err = postService.CreateCacheUser(ctx, user.ID, user.Username, user.Name)
			if err != nil {

				logger.Error("Error in inserting to cache users_cache table", zap.Error(err))

			}
			cancel()
			log.Printf("Received user registered event: ID=%v, Username=%s Name=%s,", user.ID, user.Username, user.Name)

		case err := <-pc.Errors():
			log.Printf("Consumer error encountered: %v", err)
		}
	}

}

func postLikedConsumer(consumer sarama.Consumer, postLikedTopic string, postService core.PostService, logger *platform.Logger) {

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

			if err := postService.IncreaseLikeCount(ctx, likeEvent.PostID); err != nil {

				logger.Error("Error increasing like count ", zap.Error(err))
			}

			cancel()

		case err := <-pc.Errors():

			logger.Error("post liked consumer error", zap.Error(err))
		}
	}
}

func postUnlikedConsumer(consumer sarama.Consumer, postUnLikedTopic string, postService core.PostService, logger *platform.Logger) {

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

			if err := postService.DecreaseLikeCount(ctx, likeEvent.PostID); err != nil {

				logger.Error("Error decreasing like count ", zap.Error(err))
			}

			fmt.Println("unlike coutn increased")
			cancel()

		case err := <-pc.Errors():

			logger.Error("post unliked consumer error", zap.Error(err))
		}
	}
}
