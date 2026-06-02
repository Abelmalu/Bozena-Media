package kafka

import (
	"github.com/IBM/sarama"
	"github.com/abelmalu/golang-posts/platform"
	"go.uber.org/zap"
)





func InitKafkaProducer(logger *platform.Logger){

 config := sarama.NewConfig()

 broker := []string{"localhost:9092"}

 producer, err := sarama.NewSyncProducer(broker,config)

 if err != nil {


	logger.Error("Error while initializing kafka producer",zap.Error(err))
 }
 defer producer.Close()

	logger.Info("Successfully connected to your local Kafka broker!")


	

}




