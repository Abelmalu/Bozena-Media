package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/abelmalu/golang-posts/notification/internal/core"
	"github.com/abelmalu/golang-posts/platform"
	"go.uber.org/zap"
)

type follow struct {
	ID          int `json:"id"`
	FollowerID  int `json:"follower_id"`
	FollowingID int `json:"following_id"`
}

func StartConsumer(brokers []string, userCreatedTopic, followedTopic string, notificationService core.NotificationService, logger *platform.Logger) {

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
		userCreatedConsumer(consumer, userCreatedTopic, notificationService, logger)

	}()

	go func() {

		defer wg.Done()
		followedConsumer(consumer, followedTopic, notificationService, logger)

	}()

	wg.Wait()
}

func userCreatedConsumer(consumer sarama.Consumer, userCreatedTopic string, notificationService core.NotificationService, logger *platform.Logger) {

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

			ctx, _ := context.WithTimeout(context.Background(), time.Second*2)

			err = notificationService.CreateCacheUser(ctx, user.ID, user.Username, user.Name)
			if err != nil {

				logger.Error("Error in inserting to cache users_cache table", zap.Error(err))

			}
			log.Printf("Received user registered event: ID=%v, Username=%s Name=%s,", user.ID, user.Username, user.Name)

		case err := <-pc.Errors():
			log.Printf("Consumer error encountered: %v", err)
		}

	}

}

func followedConsumer(consumer sarama.Consumer, followedTopic string, notificationService core.NotificationService, logger *platform.Logger) {

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

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

			err := notificationService.CreateUserNotification(ctx, followed.FollowerID, followed.FollowingID)

			if err != nil {

				logger.Error("Error while creating notification in the DB", zap.Error(err))
				continue

			}
			 cancel()


		case err := <-pc.Errors():

			log.Printf("Consumer error encountered: %v", err)

		}

	}

}
