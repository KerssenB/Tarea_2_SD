package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"tarea2_noti/email"
	"tarea2_noti/kafka"

	"github.com/IBM/sarama"
)

type Pedido struct {
	ID       string `json:"id"`
	Producto string `json:"producto"`
	Precio   int    `json:"precio"`
	Estado   string `json:"estado"`
}

func enviarNotificacion(pedido Pedido, smtpConfig email.SMTPConfig) {
	destinatario := "pruebaskb2024@gmail.com"
	asunto := "Actualización del Estado del Pedido"
	cuerpo := "El estado de tu pedido ha cambiado a: " + pedido.Estado

	email.EnviarCorreo(smtpConfig, destinatario, asunto, cuerpo)
}

func notificacion(smtpConfig email.SMTPConfig) {
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

			// Solo enviar notificación si el pedido tiene un estado
			if pedido.Estado != "" {
				log.Printf("Notificación - Pedido actualizado: %+v\n", pedido)
				enviarNotificacion(pedido, smtpConfig)
			}

		case <-sigchan:
			log.Println("Deteniendo el consumidor de Kafka")
			return
		}
	}
}

func main() {
	err := kafka.ConexionKafka()
	if err != nil {
		log.Fatal("Error al configurar Kafka:", err)
	}

	smtpConfig := email.SMTPConfig{
		Host:     "smtp.gmail.com",
		Port:     587,
		Username: "pruebaskb2024@gmail.com",
		Password: "sels uhky urwi rzpn",
	}

	// Iniciar el servicio de notificación
	go notificacion(smtpConfig)

	// Mantener el servicio en ejecución
	select {}
}
