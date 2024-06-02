package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
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

var (
	pedidosProcesados = make(map[string]Pedido)
	mu                sync.Mutex
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpRecipient := os.Getenv("SMTP_RECIPIENT")

	smtpPortInt, err := strconv.Atoi(smtpPort)
	if err != nil {
		log.Fatalf("Error al convertir SMTP_PORT a entero: %v", err)
	}

	config := email.NewSMTPConfig(smtpHost, smtpPortInt, smtpUsername, smtpPassword)

	err = kafka.ConexionKafka()
	if err != nil {
		log.Fatal("Error al configurar Kafka:", err)
	}

	partitions, err := kafka.Consumer.Partitions("pedidos-topic")
	if err != nil {
		log.Fatalf("Error al obtener las particiones: %v", err)
	}

	var wg sync.WaitGroup

	for _, partition := range partitions {
		wg.Add(1)
		go func(partition int32) {
			defer wg.Done()
			consumePartition(partition, config, smtpRecipient)
		}(partition)
	}

	wg.Wait()
}

func consumePartition(partition int32, config *email.SMTPConfig, recipient string) {
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

			if pedido.Estado != "" {
				mu.Lock()
				pedidosProcesados[pedido.ID] = pedido
				mu.Unlock()

				subject := "Actualización de estado del pedido"
				body := fmt.Sprintf("El pedido con ID %s ahora está en estado: %s", pedido.ID, pedido.Estado)
				err := config.SendEmail(recipient, subject, body)
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

func serveHTTP() {
	http.HandleFunc("/pedidos", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing id parameter", http.StatusBadRequest)
			return
		}

		mu.Lock()
		pedido, exists := pedidosProcesados[id]
		mu.Unlock()

		if !exists {
			http.Error(w, "Pedido not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pedido)
	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8081"
	}

	log.Printf("Servidor HTTP escuchando en el puerto %s", httpPort)
	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}
