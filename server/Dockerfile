# Указываем базовый образ, который содержит Golang
FROM golang:1.23.1-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/server.crt .
COPY --from=builder /app/server.key .

EXPOSE 8081

CMD ["./main"]