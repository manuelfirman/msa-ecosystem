# load_balancer/Dockerfile
FROM golang:1.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o load_balancer

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/load_balancer .
EXPOSE 80
CMD ["./load_balancer"]
