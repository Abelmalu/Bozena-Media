package kafka

import (
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
	"github.com/abelmalu/golang-posts/notification/internal/core"
	"github.com/abelmalu/golang-posts/platform"
)

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

			fmt.Println(msg)

		case err := <-pc.Errors():

			fmt.Println(err)

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

			fmt.Println(msg)

		case err := <-pc.Errors():

			fmt.Println(err)

		}

	}

}
