FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY user-service/go.mod user-service/go.sum ./user-service/
COPY .env ./.env

COPY user-service/ ./user-service/

WORKDIR /app/user-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/user-service/config ./config
COPY --from=builder /app/user-service/.env .
COPY --from=builder /app/.env /.env

EXPOSE 8080

CMD ["./server"]
