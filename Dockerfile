# Stage 1: Сборка
FROM golang:1.22.2 as builder
WORKDIR /http-file-server
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /myapp ./cmd/main.go

# Stage 2: Запуск
FROM alpine:latest
COPY --from=builder /myapp /myapp
EXPOSE 8082
CMD ["/myapp"]

#TRIGER PIPELINE