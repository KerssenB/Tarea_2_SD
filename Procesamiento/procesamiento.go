package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"tarea2_proce/kafka"

	"github.com/IBM/sarama"
)

type Pedido struct {
	ID       string `json:"id"`
	Producto string `json:"producto"`
	Precio   int    `json:"precio"`
	Estado   string `json:"estado"`
}

var (
	pedidosProcesados = make(map[string]Pedido)
	mu                sync.Mutex
)

func procesamiento() {
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

			log.Printf("Pedido recibido: %+v\n", pedido)

			mu.Lock()
			if _, exists := pedidosProcesados[pedido.ID]; !exists {
				pedidosProcesados[pedido.ID] = pedido
				go cambiarEstado(pedido)
			}
			mu.Unlock()

		case <-sigchan:
			log.Println("Deteniendo el consumidor de Kafka")
			return
		}
	}
}

func cambiarEstado(pedido Pedido) {
	for {
		time.Sleep(10 * time.Second) // Cambia este valor para ajustar la frecuencia

		mu.Lock()
		switch pedido.Estado {
		case "":
			pedido.Estado = "Recibido"
		case "Recibido":
			pedido.Estado = "Preparando"
		case "Preparando":
			pedido.Estado = "Entregado"
		case "Entregado":
			pedido.Estado = "Finalizado"
		default:
			mu.Unlock()
			return // Si el estado es "finalizado", terminar la función
		}

		// Actualizar el pedido en el mapa antes de enviar
		pedidosProcesados[pedido.ID] = pedido
		mu.Unlock()

		// Verificar el estado antes de enviar
		if pedido.Estado == "" {
			continue
		}

		pedidoJSON, err := json.Marshal(pedido)
		if err != nil {
			log.Println("Error al convertir el pedido a JSON:", err)
			continue
		}

		producer := kafka.Producer
		_, _, err = producer.SendMessage(&sarama.ProducerMessage{
			Topic: "pedidos-topic",
			Value: sarama.StringEncoder(string(pedidoJSON)),
		})
		if err != nil {
			log.Println("Error al enviar el mensaje a Kafka:", err)
			continue
		}

		log.Printf("Pedido actualizado y enviado a Kafka: %+v\n", pedido)
	}
}

func main() {
	err := kafka.ConexionKafka()
	if err != nil {
		log.Fatal("Error al configurar Kafka:", err)
	}

	// Iniciar el servicio de procesamiento
	go procesamiento()

	// Mantener el servicio en ejecución
	select {}
}
