package kafka

import (
	"github.com/IBM/sarama"
	 "github.com/abelmalu/golang-posts/like/internal/errors"
	"github.com/abelmalu/golang-posts/platform"
	"go.uber.org/zap"
)

func InitKafkaProducer(logger *platform.Logger) (sarama.SyncProducer, error) {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	broker := []string{"localhost:9092"}

	producer, err := sarama.NewSyncProducer(broker, config)

	if err != nil {

		logger.Error("Error while initializing kafka producer", zap.Error(err))

		return nil, ierrors.ErrKafkaConnection
	}

	logger.Info("Kafka connected successfully!")

	return producer, nil

}
