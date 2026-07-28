# product-service

Product Catalog Service: manages product listings, categories, and inventory.

## Stack

- Go, Gin (HTTP), PostgreSQL (sqlx), viper (config), log/slog (JSON logs)

## Structure

- `cmd/` — entrypoint (`main.go`) and app wiring (`app/app.go`)
- `internal/config/` — configuration (yaml + `.env` secrets)
- `internal/transport/http/` — HTTP server, `Handler` interface
- `pkg/database/` — PostgreSQL connection
- `migrations/` — goose migrations (categories, products, inventory)

## Run locally

```bash
cp .env.example .env   # set postgres_password
docker compose up
```

HTTP: `:8082` (`GET /health`), database: `hookah_products` on `:5434`.
