include .env
export

# Docker
up:
	docker-compose up -d

down:
	docker-compose down

# Database migrations
migrate-up:
	goose -dir=$(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	goose -dir=$(MIGRATIONS_DIR) postgres "$(DB_URL)" down

# Lint
linter:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.2 run

# Start
run:
	go run main.go


