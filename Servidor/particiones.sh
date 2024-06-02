while ! nc -z kafka 9092; do   
  sleep 1
done

kafka-topics.sh --create --topic pedidos-topic --bootstrap-server kafka:9092 --partitions 3 --replication-factor 1 --if-not-exists
