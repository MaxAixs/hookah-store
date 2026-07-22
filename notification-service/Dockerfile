FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY notification-service/go.mod notification-service/go.sum ./notification-service/
COPY user-service/go.mod user-service/go.sum ./user-service/

COPY notification-service/ ./notification-service/
COPY user-service/pkg/ ./user-service/pkg/
COPY .env ./.env

WORKDIR /app/notification-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/main.go

FROM alpine:3.20
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/notification-service/config ./config
COPY --from=builder /app/notification-service/.env /app/.env
COPY --from=builder /app/.env /.env
EXPOSE 8081
CMD ["./server"]
