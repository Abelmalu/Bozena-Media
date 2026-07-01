package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/abelmalu/golang-posts/Auth/internal/core"
	"github.com/abelmalu/golang-posts/platform"
	"go.uber.org/zap"
)

type follow struct {
	ID          int `json:"id"`
	FollowerID  int `json:"follower_id"`
	FollowingID int `json:"following_id"`
}

func StartConsumer(brokers []string, followedTopic, unfollowedTopic string, authService core.AuthService, logger *platform.Logger) {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating master consumer: %v", err)
	}

	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {

		defer wg.Done()

		followedConsumer(consumer, followedTopic, authService, logger)
	}()

	go func() {

		defer wg.Done()

		unfollowedConsumer(consumer, unfollowedTopic, authService, logger)

	}()
}

func followedConsumer(consumer sarama.Consumer, followedTopic string, authService core.AuthService, logger *platform.Logger) {

	pc, err := consumer.ConsumePartition(followedTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error consuming post partition: %v", err)
	}
	defer pc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for {

		select {

		case msg := <-pc.Messages():

			var followed follow

			if err := json.Unmarshal(msg.Value, &followed); err != nil {

				logger.Error("Error while unmarshaling followed event", zap.Error(err))
				continue
			}

			if err := authService.IncreaseFollowCounts(ctx, followed.FollowerID, followed.FollowingID); err != nil {

				logger.Error("Error while updating user followed event", zap.Error(err))
			}

		case err := <-pc.Errors():
			logger.Error("Post consumer error: %v", zap.Error(err))
		}
	}

}

func unfollowedConsumer(consumer sarama.Consumer, unfollowedTopic string, authService core.AuthService, logger *platform.Logger) {

	pc, err := consumer.ConsumePartition(unfollowedTopic, 0, sarama.OffsetNewest)
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

		case err := <-pc.Errors():
			logger.Error("Post consumer error: %v", zap.Error(err))
		}
	}

}
