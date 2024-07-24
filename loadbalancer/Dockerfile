# load_balancer/Dockerfile
FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o load_balancer

FROM debian:latest
WORKDIR /root/
COPY --from=builder /app/load_balancer .
RUN chmod +x /root/load_balancer
EXPOSE 80
CMD ["./load_balancer"]
