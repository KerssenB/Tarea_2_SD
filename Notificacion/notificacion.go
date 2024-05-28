package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"tarea2_noti/email"
	"tarea2_noti/kafka"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

type Pedido struct {
	ID       string `json:"id"`
	Producto string `json:"producto"`
	Precio   int    `json:"precio"`
	Estado   string `json:"estado"`
}

func main() {
	// Cargar las variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	// Obtener la configuración SMTP desde las variables de entorno
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpRecipient := os.Getenv("SMTP_RECIPIENT")

	// Convertir el puerto SMTP a entero
	smtpPortInt, err := strconv.Atoi(smtpPort)
	if err != nil {
		log.Fatalf("Error al convertir SMTP_PORT a entero: %v", err)
	}

	// Configurar SMTP
	config := email.NewSMTPConfig(smtpHost, smtpPortInt, smtpUsername, smtpPassword)

	err = kafka.ConexionKafka()
	if err != nil {
		log.Fatal("Error al configurar Kafka:", err)
	}

	partitionConsumer, err := kafka.Consumer.ConsumePartition("pedidos-topic", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error al iniciar el consumidor de la partición: %v", err)
	}
	defer partitionConsumer.Close()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var pedido Pedido
			err := json.Unmarshal(msg.Value, &pedido)
			if err != nil {
				log.Printf("Error al deserializar el mensaje: %v", err)
				continue
			}

			if pedido.Estado != "" {
				subject := "Actualización de estado del pedido"
				body := fmt.Sprintf("El pedido con ID %s ahora está en estado: %s", pedido.ID, pedido.Estado)
				err := config.SendEmail(smtpRecipient, subject, body)
				if err != nil {
					log.Printf("Error al enviar el correo: %v", err)
				} else {
					log.Printf("Notificación - Pedido actualizado: %+v\n", pedido)
				}
			}

		case <-sigchan:
			log.Println("Deteniendo el consumidor de Kafka")
			return
		}
	}
}
