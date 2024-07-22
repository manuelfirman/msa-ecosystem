# Etapa de construcción
FROM golang:1.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o auth_service

# Etapa de ejecución
FROM ubuntu:latest
WORKDIR /root/
COPY --from=builder /app/auth_service /root/auth_service
RUN ls -l /root/ # Verifica que el archivo está ahí
RUN chmod +x /root/auth_service # Asegúrate de que el archivo sea ejecutable
EXPOSE 5000
CMD ["/root/auth_service"]
