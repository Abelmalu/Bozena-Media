package kafka

import (
	"github.com/IBM/sarama"
	ierrors "github.com/abelmalu/golang-posts/Auth/internal/errors"
	"github.com/abelmalu/golang-posts/platform"
	"go.uber.org/zap"
)





func InitKafkaProducer(logger *platform.Logger) ( sarama.SyncProducer,error ){

 config := sarama.NewConfig()

 broker := []string{"localhost:9092"}

 producer, err := sarama.NewSyncProducer(broker,config)

 if err != nil {


	logger.Error("Error while initializing kafka producer",zap.Error(err))

	return nil, ierrors.ErrKafkaConnection
 }
 defer producer.Close()

	logger.Info("Successfully connected to your local Kafka broker!")


	return producer,nil

}




