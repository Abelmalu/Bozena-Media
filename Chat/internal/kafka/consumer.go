package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/abelmalu/golang-posts/Chat/internal/core"
	"github.com/abelmalu/golang-posts/platform"
	"go.uber.org/zap"
)

func StartConsumer(brokers []string, userCreatedtopic string, cs core.ChatService, logger *platform.Logger) {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating master consumer: %v", err)
	}

	go func() {

		userCreatedConsumer(consumer, userCreatedtopic, cs, logger)
	}()

}

func userCreatedConsumer(consumer sarama.Consumer, userCreatedTopic string, cs core.ChatService, logger *platform.Logger) {

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
				Avatar   string `json:"avatar"`
			}

			err := json.Unmarshal(msg.Value, &user)
			if err != nil {
				log.Printf("Failed to unmarshal user data: %v", err)
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

			err = cs.CreateCacheUser(ctx, user.ID, user.Username, user.Name, user.Avatar)
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
