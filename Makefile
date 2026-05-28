include .env
export

migrate-up:
	migrate -path internal/database/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" up

migrate-down:
	migrate -path internal/database/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" down

# ── Seeders ────────────────────────────────────────────────────
seed:
	go run cmd/seed/main.go

seed-fresh:
	go run cmd/seed/main.go -fresh

seed-clean:
	go run cmd/seed/main.go -clean

seed-admin:
	go run cmd/seed/main.go -admin

# ── Development ────────────────────────────────────────────────
dev:
	go run cmd/api/main.go
