.PHONY: dev build run test clean docker-up docker-down migrate

# Development
dev:
	go run ./cmd/server

# Build
build:
	go build -o bin/nyengo ./cmd/server

# Run built binary
run: build
	./bin/nyengo

# Test
test:
	go test -v ./...

# Clean
clean:
	rm -rf bin/

# Docker
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build

# Database
migrate-up:
	migrate -path internal/database/migrations -database "$$DATABASE_URL" up

migrate-down:
	migrate -path internal/database/migrations -database "$$DATABASE_URL" down

# Swagger
swagger:
	swag init -g cmd/server/main.go -o api/docs

# Tidy
tidy:
	go mod tidy
