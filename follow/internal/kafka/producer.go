package kafka

import (
	"github.com/IBM/sarama"
	ierrors "github.com/abelmalu/golang-posts/follow/internal/errors"
	"github.com/abelmalu/golang-posts/follow/config"
	"github.com/abelmalu/golang-posts/platform"
	"go.uber.org/zap"
)


func InitKafkaProducer(logger *platform.Logger,cfg *config.Config) ( sarama.SyncProducer,error ){

 config := sarama.NewConfig()
  config.Producer.Return.Successes = true

 producer, err := sarama.NewSyncProducer(cfg.KafkaBrokersURL,config)

 if err != nil {


	logger.Error("Error while initializing kafka producer",zap.Error(err))

	return nil, ierrors.ErrKafkaConnection
 }

	logger.Info("Kafka connected successfully!")


	return producer,nil

}