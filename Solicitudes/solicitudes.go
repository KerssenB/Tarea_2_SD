package main

import (
	"encoding/json"
	"log"
	"net/http"

	"tarea2_sol/kafka"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Pedido struct {
	ID       string `json:"id"`
	Producto string `json:"producto"`
	Precio   int    `json:"precio"`
}

func crearPedido(w http.ResponseWriter, r *http.Request) {
	var nuevoPedido Pedido
	err := json.NewDecoder(r.Body).Decode(&nuevoPedido)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nuevoPedido.ID = uuid.New().String()

	log.Printf("Pedido recibido: %+v\n", nuevoPedido)

	pedidoJSON, err := json.Marshal(nuevoPedido)
	if err != nil {
		log.Println("Error al convertir el pedido a JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	producer := kafka.Producer
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "pedidos-topic",
		Value: sarama.StringEncoder(string(pedidoJSON)),
	})
	if err != nil {
		log.Println("Error al enviar el mensaje a Kafka:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Pedido recibido correctamente", "id": nuevoPedido.ID})
}

func main() {
	err := kafka.ConexionKafka()
	if err != nil {
		log.Fatal("Error al configurar Kafka:", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/pedidos", crearPedido).Methods("POST")

	log.Println("Servidor iniciado en el puerto 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
