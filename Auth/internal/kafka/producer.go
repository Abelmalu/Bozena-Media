package kafka

import (
	"github.com/IBM/sarama"
	ierrors "github.com/abelmalu/golang-posts/Auth/internal/errors"
	"github.com/abelmalu/golang-posts/platform"
	"go.uber.org/zap"
)

func InitKafkaProducer(logger *platform.Logger, brokers []string) (sarama.SyncProducer, error) {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)

	if err != nil {

		logger.Error("Error while initializing kafka producer", zap.Error(err))

		return nil, ierrors.ErrKafkaConnection
	}

	logger.Info("Kafka connected successfully!")

	return producer, nil

}
