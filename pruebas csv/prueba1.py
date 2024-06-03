import requests
import json
import time
import csv

# Configuración
base_url = "http://localhost:8080/pedidos"  # URL base del servidor HTTP
intervalo = 0.1  # Intervalo de tiempo entre cada pedido en segundos
output_file = "test_carga_resultados.csv"  # Nombre del archivo CSV para guardar los resultados

# Función para enviar pedidos en intervalos regulares y guardar resultados
def enviar_pedidos():
    with open(output_file, mode='w', newline='') as file:
        writer = csv.writer(file)
        writer.writerow(["Pedido", "Estado", "Tiempo de respuesta", "Código de estado", "Tamaño de la solicitud (bytes)", "Tamaño de la respuesta (bytes)", "Timestamp de inicio", "Timestamp de fin"])

        # Carga de datos desde el archivo JSON
        with open("Datasets/productos.json", "r") as f:
            data = json.load(f)

        # Iteración sobre cada producto en el archivo JSON
        for i, pedido in enumerate(data):
            start_time = time.time()
            request_size = len(json.dumps(pedido))
            response = requests.post(base_url, json=pedido)
            end_time = time.time()
            response_size = len(response.content)

            tiempo_respuesta = end_time - start_time
            estado = "Creado correctamente" if response.status_code == 201 else "Error al crear"
            
            writer.writerow([i+1, estado, tiempo_respuesta, response.status_code, request_size, response_size, start_time, end_time])
            time.sleep(intervalo)

# Envío de pedidos y guardado de resultados en CSV
enviar_pedidos()