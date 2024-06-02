#!/bin/bash

# Esperar a que Kafka esté disponible
while ! nc -z kafka 9092; do   
  sleep 1
done

# Crear el topic con múltiples particiones
kafka-topics.sh --create --topic pedidos-topic --bootstrap-server kafka:9092 --partitions 3 --replication-factor 1 --if-not-exists
