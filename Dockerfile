# Etapa de construcción
FROM golang:1.23.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /go-api ./cmd

# Etapa de ejecución
FROM alpine:latest
WORKDIR /
COPY --from=builder /go-api /go-api
EXPOSE 8080
CMD ["/go-api"]
