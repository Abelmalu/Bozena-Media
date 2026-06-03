package main

import (
	//"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

// Define your User struct to unmarshal the incoming bytes back into Go data
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	// Add other fields matching your user schema
}

// StartConsumer initializes and runs a simple Kafka partition consumer
func StartConsumer(brokers []string, topic string) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

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
			var user User
			
			// Unmarshal the raw JSON byte string back into your Go struct
			err := json.Unmarshal(msg.Value, &user)
			if err != nil {
				log.Printf("Failed to unmarshal user data: %v", err)
				continue
			}

			// Now you can use the object natively in your code
			log.Printf("Received user registered event: ID=%s, Name=%s", user.ID, user.Name)

		case err := <-partitionConsumer.Errors():
			log.Printf("Consumer error encountered: %v", err)
		}
	}
}
