package kafka


import (
	"encoding/json"
	"log"
	"github.com/IBM/sarama"

)

func StartConsumer(brokers []string, topic string) {
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
			}
			
			// Unmarshal the raw JSON byte string back into your Go struct
			err := json.Unmarshal(msg.Value, &user)
			if err != nil {
				log.Printf("Failed to unmarshal user data: %v", err)
				continue
			}

			// Now you can use the object natively in your code
			log.Printf("Received user registered event: ID=%v, Name=%s", user.ID, user.Username)

		case err := <-partitionConsumer.Errors():
			log.Printf("Consumer error encountered: %v", err)
		}
	}
}