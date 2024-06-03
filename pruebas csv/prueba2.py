import requests
import json
import time
import threading
import csv

base_url = "http://localhost:8080/pedidos" 
num_pedidos = 1000  
output_file = "test_gran_cantidad_resultados.csv" 

def cargar_datos_productos():
   
    with open("Datasets/productos.json", "r") as f:
        data = json.load(f)
    return data

def enviar_lote_pedidos(inicio, fin, writer, data):
    productos = data[inicio-1:fin]
    for i, pedido in enumerate(productos, inicio):
        start_time = time.time()
        request_size = len(json.dumps(pedido))
        response = requests.post(base_url, json=pedido)
        end_time = time.time()
        response_size = len(response.content)

        tiempo_respuesta = end_time - start_time
        estado = "Creado correctamente" if response.status_code == 201 else "Error al crear"
        
        writer.writerow([i, estado, tiempo_respuesta, response.status_code, request_size, response_size, start_time, end_time])

data = cargar_datos_productos()

with open(output_file, mode='w', newline='') as file:
    writer = csv.writer(file)
    writer.writerow(["Pedido", "Estado", "Tiempo de respuesta", "Código de estado", "Tamaño de la solicitud (bytes)", "Tamaño de la respuesta (bytes)", "Timestamp de inicio", "Timestamp de fin"])

    threads = [] 
    for i in range(0, num_pedidos, 100):
        t = threading.Thread(target=enviar_lote_pedidos, args=(i+1, i+101, writer, data))
        threads.append(t)
        t.start()

    for t in threads:
        t.join()