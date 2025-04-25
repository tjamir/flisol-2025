FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# Copia tudo do contexto (microblog/)
COPY . .

# Entra no timeline-worker para rodar os comandos
WORKDIR /app/timeline-worker

RUN go mod download
RUN go build -o /app/app ./cmd

FROM alpine:latest

WORKDIR /root/
RUN apk add --no-cache bash
COPY --from=builder /app/app .
COPY --from=builder /app/timeline-worker/wait-for-it.sh .

EXPOSE 8085

CMD ["./wait-for-it.sh", "kafka:9092", "--", "./wait-for-it.sh", "cassandra:9042", "--", "./wait-for-it.sh", "follow-service:8083", "--", "./wait-for-it.sh", "post-service:8082", "--", "./app"]
