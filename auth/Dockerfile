FROM golang:1.18 AS builder
WORKDIR /app
COPY . .

WORKDIR /app/cmd
RUN go build -o /app/auth_service

# Etapa de ejecución
FROM debian:latest
WORKDIR /root/

# Copiar el binario de la etapa de construcción
COPY --from=builder /app/auth_service /root/auth_service

# Copiar el script de configuración
COPY env.sh /root/env.sh

# Hacer el binario ejecutable
RUN chmod +x /root/auth_service

# Asegurarse de que el script de configuración también sea ejecutable
RUN chmod +x /root/env.sh

# Exponer el puerto
EXPOSE 5000

# Configurar el punto de entrada para ejecutar el script de configuración y luego la aplicación
CMD ["/bin/sh", "-c", ". /root/env.sh && /root/auth_service"]
