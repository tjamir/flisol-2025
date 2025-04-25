# Etapa de build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Instala git para permitir o uso de módulos Go
RUN apk add --no-cache git

# Copia dependências e baixa os módulos
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copia o restante dos arquivos e compila
COPY . .
RUN go build -o app ./cmd

# Etapa de execução
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/app .

# Exponha a porta correta no docker-compose (se necessário)
EXPOSE 8080

CMD ["./app"]
