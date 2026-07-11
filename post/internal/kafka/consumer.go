package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/abelmalu/golang-posts/post/internal/core"
	"go.uber.org/zap"
)

func StartConsumer(brokers []string, topic string, postService core.PostService,logger *platform.Logger) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	log.Printf("Consumer started. Listening on topic: %s...", topic)

	for {

		select {
		case msg := <-partitionConsumer.Messages():
			var user struct {

				ID int `json:"id"`
				Username string `json:"username"`
				Name string `json:"name"`
			}
			
			err := json.Unmarshal(msg.Value, &user)
			if err != nil {
				log.Printf("Failed to unmarshal user data: %v", err)
				continue
			}

			
			ctx,_ := context.WithTimeout(context.Background(),time.Second*2)

			err = postService.CreateCacheUser(ctx,user.ID,user.Username,user.Name)
			if err != nil {

				logger.Error("Error in inserting to cache users_cache table",zap.Error(err))

				
			}
			log.Printf("Received user registered event: ID=%v, Username=%s Name=%s,", user.ID, user.Username,user.Name)

		case err := <-partitionConsumer.Errors():
			log.Printf("Consumer error encountered: %v", err)
		}
	}  
}