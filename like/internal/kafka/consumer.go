package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/abelmalu/golang-posts/like/internal/core"
	"go.uber.org/zap"
)

func StartConsumer(brokers []string, topic string, likeService core.LikeService,logger *platform.Logger) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// 1. Create the master consumer
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	defer consumer.Close()

	// 2. Consume from partition 0 (use a ConsumerGroup if scaling horizontally)
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	log.Printf("Consumer started. Listening on topic: %s...", topic)

	// 3. Process loop
	for {

		select {
		case msg := <-partitionConsumer.Messages():
			var user struct {

				ID int `json:"id"`
				Username string `json:"username"`
				Name string `json:"name"`
			}
			
			// Unmarshal the raw JSON byte string back into your Go struct
			err := json.Unmarshal(msg.Value, &user)
			if err != nil {
				log.Printf("Failed to unmarshal user data: %v", err)
				continue
			}

			
			ctx,cancel := context.WithTimeout(context.Background(),time.Second*2)
			defer cancel()

			err = likeService.CreateCacheUser(ctx,user.ID,user.Username,user.Name)
			if err != nil {

				logger.Error("Error in inserting to cache users_cache table",zap.Error(err))

				
			}
			log.Printf("Received user registered event: ID=%v, Username=%s Name=%s,", user.ID, user.Username,user.Name)

		case err := <-partitionConsumer.Errors():
			log.Printf("Consumer error encountered: %v", err)
		}
	}  
}