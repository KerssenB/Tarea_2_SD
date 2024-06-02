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

func cambiarEstado(pedido Pedido) {
	for {
		time.Sleep(10 * time.Second)

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
			return
		}

		pedidosProcesados[pedido.ID] = pedido
		mu.Unlock()

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
	}
}

func procesamiento() {
	partitions, err := kafka.Consumer.Partitions("pedidos-topic")
	if err != nil {
		log.Fatalf("Error al obtener las particiones: %v", err)
	}

	var wg sync.WaitGroup

	for _, partition := range partitions {
		wg.Add(1)
		go func(partition int32) {
			defer wg.Done()
			consumePartition(partition)
		}(partition)
	}

	wg.Wait()
}

func consumePartition(partition int32) {
	partitionConsumer, err := kafka.Consumer.ConsumePartition("pedidos-topic", partition, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error al iniciar el consumidor de la partición %d: %v", partition, err)
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

			log.Printf("Pedido recibido de la partición %d: %+v\n", partition, pedido)

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

func main() {
	err := kafka.ConexionKafka()
	if err != nil {
		log.Fatal("Error al configurar Kafka:", err)
	}

	go procesamiento()

	select {}
}
