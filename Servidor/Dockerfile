FROM bitnami/kafka:latest

COPY particiones.sh /usr/local/bin/particiones.sh

CMD ["/bin/bash", "-c", "chmod +x /usr/local/bin/particiones.sh && /usr/local/bin/particiones.sh & /opt/bitnami/scripts/kafka/run.sh"]
