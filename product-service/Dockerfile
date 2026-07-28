FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY product-service/go.mod product-service/go.sum ./product-service/
COPY .env ./.env

COPY product-service/ ./product-service/

WORKDIR /app/product-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/product-service/config ./config
COPY --from=builder /app/product-service/.env .
COPY --from=builder /app/.env /.env

EXPOSE 8082

CMD ["./server"]
