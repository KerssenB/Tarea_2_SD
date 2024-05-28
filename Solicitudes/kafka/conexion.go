package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

var Producer sarama.SyncProducer

func ConexionKafka() error {
	config := sarama.NewConfig()
	config.ClientID = "mi-cliente"
	config.Producer.Return.Successes = true

	// Utilizar el nombre del servicio de Kafka en la red de Docker
	brokers := []string{"kafka:9092"}

	var err error
	Producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return err
	}

	log.Println("Conexión a Kafka exitosa")
	return nil
}
