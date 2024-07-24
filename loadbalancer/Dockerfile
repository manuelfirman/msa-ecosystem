# load_balancer/Dockerfile
FROM golang:1.22 AS builder
WORKDIR /app
COPY . .

WORKDIR /app/cmd
RUN go build -o /app/load_balancer

FROM debian:latest
WORKDIR /root/

COPY env.sh /root/env.sh
RUN chmod +x /root/env.sh

COPY --from=builder /app/load_balancer /root/load_balancer
RUN chmod +x /root/load_balancer

EXPOSE 80

CMD ["/bin/sh", "-c", ". /root/env.sh && /root/load_balancer"]
