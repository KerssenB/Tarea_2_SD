import random
import json
import itertools

# Lista de nombres de productos en español
nombres_productos = [
    "Camiseta",
    "Pantalones",
    "Zapatos",
    "Bolso",
    "Sombrero",
    "Bufanda",
    "Gafas de sol",
    "Reloj",
    "Pendientes",
    "Collar",
    "Anillo",
    "Cinturón",
    "Chaqueta",
    "Vestido",
    "Falda",
    "Calcetines",
    "Guantes",
    "Pijama",
    "Abrigo",
    "Sudadera"
]

# Lista de colores
colores = ["Negro", "Blanco", "Azul", "Rojo", "Verde", "Gris", "Amarillo", "Morado", "Rosado", "Naranja"]

# Diccionario de rangos de precios por categoría de productos en CLP
rangos_precios = {
    "Camiseta": (5000, 20000),
    "Pantalones": (10000, 30000),
    "Zapatos": (20000, 60000),
    "Bolso": (10000, 50000),
    "Sombrero": (5000, 20000),
    "Bufanda": (5000, 20000),
    "Gafas de sol": (10000, 40000),
    "Reloj": (20000, 80000),
    "Pendientes": (5000, 30000),
    "Collar": (10000, 50000),
    "Anillo": (5000, 20000),
    "Cinturón": (5000, 20000),
    "Chaqueta": (20000, 80000),
    "Vestido": (20000, 80000),
    "Falda": (10000, 40000),
    "Calcetines": (1000, 5000),
    "Guantes": (5000, 20000),
    "Pijama": (10000, 40000),
    "Abrigo": (30000, 100000),
    "Sudadera": (10000, 40000)
}

# Generar todas las combinaciones posibles de productos con sus precios
combinaciones = []
for producto, color in itertools.product(nombres_productos, colores):
    precio_min, precio_max = rangos_precios[producto]
    precio = random.randint(precio_min, precio_max) 
    combinaciones.append({"producto": f"{producto} {color}", "precio": precio})

# Seleccionar 300,000 productos de forma aleatoria
productos = []
total_c = len(combinaciones) - 1

for i in range(300000):
    n = random.randint(0, total_c)
    productos.append(combinaciones[n])

# Escribir los productos seleccionados en un archivo JSON
with open('productos.json', mode='w', encoding='utf-8') as file:
    json.dump(productos, file, ensure_ascii=False, indent=4)

print("Se han seleccionado y guardado 300,000 productos en 'productos.json'.")
