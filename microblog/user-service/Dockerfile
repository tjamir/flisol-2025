FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# Copia tudo do contexto (microblog/)
COPY . .

# Entra no user-service para rodar os comandos
WORKDIR /app/user-service

RUN go mod download
RUN go build -o /app/app ./cmd

FROM alpine:latest

WORKDIR /root/
RUN apk add --no-cache bash
COPY --from=builder /app/app .
COPY --from=builder /app/user-service/wait-for-it.sh .

EXPOSE 8081

CMD ["./wait-for-it.sh", "db:5432", "--", "./app"]
