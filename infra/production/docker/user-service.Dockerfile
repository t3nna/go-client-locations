FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
WORKDIR /app/services/user-service
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/services/user-service/user-service .
CMD ["./user-service"]