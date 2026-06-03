package main

import (
	"log"

	"github.com/IBM/sarama"
)

func main() {

	InitKafkaProducer()

}

func InitKafkaProducer() {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	broker := []string{"localhost:9092"}

	producer, err := sarama.NewSyncProducer(broker, config)

	if err != nil {

		log.Fatal("Error while initializing kafka producer", err)
	}
	defer producer.Close()

	log.Println("Successfully connected to your local Kafka broker!")

	msg := &sarama.ProducerMessage{
		Topic: "userCreated",
		Value: sarama.StringEncoder("Hello, Kafka from Go!"),
	}

	// 5. Send the message to Kafka
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	log.Printf("Message successfully sent!\n")
	log.Printf("Topic: %s | Partition: %d | Offset: %d\n", msg.Topic, partition, offset)

}
