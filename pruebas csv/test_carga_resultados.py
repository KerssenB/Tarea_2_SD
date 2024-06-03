import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns

# Leer el archivo CSV
df = pd.read_csv('test_carga_resultados.csv')

# Calcular estadísticas descriptivas
media_tiempo_respuesta = df['Tiempo de respuesta'].mean()
mediana_tiempo_respuesta = df['Tiempo de respuesta'].median()
desviacion_tiempo_respuesta = df['Tiempo de respuesta'].std()
rango_intercuartil_tiempo_respuesta = df['Tiempo de respuesta'].quantile(0.75) - df['Tiempo de respuesta'].quantile(0.25)

media_tamaño_solicitud = df['Tamaño de la solicitud (bytes)'].mean()
mediana_tamaño_solicitud = df['Tamaño de la solicitud (bytes)'].median()
desviacion_tamaño_solicitud = df['Tamaño de la solicitud (bytes)'].std()
rango_intercuartil_tamaño_solicitud = df['Tamaño de la solicitud (bytes)'].quantile(0.75) - df['Tamaño de la solicitud (bytes)'].quantile(0.25)

media_tamaño_respuesta = df['Tamaño de la respuesta (bytes)'].mean()
mediana_tamaño_respuesta = df['Tamaño de la respuesta (bytes)'].median()
desviacion_tamaño_respuesta = df['Tamaño de la respuesta (bytes)'].std()
rango_intercuartil_tamaño_respuesta = df['Tamaño de la respuesta (bytes)'].quantile(0.75) - df['Tamaño de la respuesta (bytes)'].quantile(0.25)

# Imprimir estadísticas
print("Media del tiempo de respuesta:", media_tiempo_respuesta)
print("Mediana del tiempo de respuesta:", mediana_tiempo_respuesta)
print("Desviación estándar del tiempo de respuesta:", desviacion_tiempo_respuesta)
print("Rango intercuartil del tiempo de respuesta:", rango_intercuartil_tiempo_respuesta)

print("Media del tamaño de solicitud:", media_tamaño_solicitud)
print("Mediana del tamaño de solicitud:", mediana_tamaño_solicitud)
print("Desviación estándar del tamaño de solicitud:", desviacion_tamaño_solicitud)
print("Rango intercuartil del tamaño de solicitud:", rango_intercuartil_tamaño_solicitud)

print("Media del tamaño de respuesta:", media_tamaño_respuesta)
print("Mediana del tamaño de respuesta:", mediana_tamaño_respuesta)
print("Desviación estándar del tamaño de respuesta:", desviacion_tamaño_respuesta)
print("Rango intercuartil del tamaño de respuesta:", rango_intercuartil_tamaño_respuesta)

# Histogramas
plt.figure(figsize=(15, 5))

plt.subplot(1, 3, 1)
sns.histplot(df['Tiempo de respuesta'], bins=30)
plt.title('Histograma del tiempo de respuesta')

plt.subplot(1, 3, 2)
sns.histplot(df['Tamaño de la solicitud (bytes)'], bins=30)
plt.title('Histograma del tamaño de solicitud')

plt.subplot(1, 3, 3)
sns.histplot(df['Tamaño de la respuesta (bytes)'], bins=30)
plt.title('Histograma del tamaño de respuesta')

plt.tight_layout()
plt.show()

# Boxplots
plt.figure(figsize=(15, 5))

plt.subplot(1, 3, 1)
sns.boxplot(y=df['Tiempo de respuesta'])
plt.title('Boxplot del tiempo de respuesta')

plt.subplot(1, 3, 2)
sns.boxplot(y=df['Tamaño de la solicitud (bytes)'])
plt.title('Boxplot del tamaño de solicitud')

plt.subplot(1, 3, 3)
sns.boxplot(y=df['Tamaño de la respuesta (bytes)'])
plt.title('Boxplot del tamaño de respuesta')

plt.tight_layout()
plt.show()

# Diagramas de puntos y análisis multivariado
plt.figure(figsize=(15, 5))

plt.subplot(1, 2, 1)
sns.scatterplot(x=df['Tamaño de la solicitud (bytes)'], y=df['Tiempo de respuesta'], hue=df['Estado'] == "Creado correctamente")
plt.title('Tamaño de solicitud vs Tiempo de respuesta')

plt.subplot(1, 2, 2)
sns.scatterplot(x=df['Tamaño de la solicitud (bytes)'], y=df['Tiempo de respuesta'])
plt.title('Tamaño de solicitud vs Tiempo de respuesta')

plt.tight_layout()
plt.show()

# Correlaciones
correlacion_creado_correctamente = df[df['Estado'] == "Creado correctamente"]['Tamaño de la solicitud (bytes)'].corr(df['Tiempo de respuesta'])
correlacion_general = df['Tamaño de la solicitud (bytes)'].corr(df['Tiempo de respuesta'])

print("Índice de correlación entre tamaño de solicitud y tiempo de respuesta (creado correctamente):", correlacion_creado_correctamente)
print("Índice de correlación entre tamaño de solicitud y tiempo de respuesta (general):", correlacion_general)

# Matriz de correlación
correlation_matrix = df[['Tiempo de respuesta', 'Tamaño de la solicitud (bytes)', 'Tamaño de la respuesta (bytes)']].corr()
print("Matriz de correlación:")
print(correlation_matrix)

# Ancho de violín
plt.figure(figsize=(15, 5))

plt.subplot(1, 3, 1)
sns.violinplot(y=df['Tiempo de respuesta'])
plt.title('Ancho de violín del tiempo de respuesta')

plt.subplot(1, 3, 2)
sns.violinplot(y=df['Tamaño de la solicitud (bytes)'])
plt.title('Ancho de violín del tamaño de solicitud')

plt.subplot(1, 3, 3)
sns.violinplot(y=df['Tamaño de la respuesta (bytes)'])
plt.title('Ancho de violín del tamaño de respuesta')

plt.tight_layout()
plt.show()
