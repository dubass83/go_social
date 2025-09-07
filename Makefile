.PHONY: *

# include .env.cloud
include .env
export
BINARY_NAME=myapp

## build: Build binary
build:
	@echo "Building..."
	env CGO_ENABLED=0  go build -ldflags="-s -w" -o bin/${BINARY_NAME} ./cmd/api
	@echo "Built!"

## run: go run
run: start_db migrate_up seed
	@echo "Starting..."
	go run ./cmd/api

## init_run: go run with initial Seeding the database
init_run: start_db migrate_up seed
	@echo "Starting..."
	go run ./cmd/api

## clean: runs go clean and deletes binaries
clean: stop_db
	@echo "Cleaning..."
	@go clean
	@rm bin/${BINARY_NAME} || true
	@rm bin/main || true
	@rm bin/build-errors.log || true
	@echo "Cleaned!"

## start: build and run compiled app
start: build start_db migrate_up
	@echo "Starting..."
	@env ./bin/${BINARY_NAME} &
	@echo "Started!"

## stop: stops the running application
stop: stop_db
	@echo "Stopping..."
	@-pkill -SIGTERM -f "./bin/${BINARY_NAME}"
	@echo "Stopped!"

## restart: stops and starts the application
restart: stop start

## repl: starts the application in REPL mode
repl: start_db migrate_up
	air

## test: runs all tests
test: start_db migrate_up
	go test -v ./...

race_test: start_db
	go test -race -v ./...

start_db:
	@echo "Starting database..."
	@docker-compose up -d db
	while ! docker-compose exec db pg_isready -h localhost -p 5432; do sleep 1; done
	@docker-compose exec db psql -U postgres -c "CREATE DATABASE social;" 2>/dev/null || echo "Database might already exist"
	@echo "Database started!"

stop_db:
	@echo "Stopping database..."
	@docker-compose down -v db
	@echo "Database stopped!"

start_redis:
	@echo "Starting Redis..."
	@docker-compose up -d redis
	@echo "Redis started!"

stop_redis:
	@echo "Stopping Redis..."
	@docker-compose down -v redis
	@echo "Redis stopped!"

new_migration:
	migrate create -ext sql -dir cmd/db/migration -seq ${name}

migrate_up:
	migrate -path cmd/db/migration -database ${DB_ADDR} -verbose up

migrate_up1:
	migrate -path cmd/db/migration -database ${DB_ADDR} -verbose up 1

migrate_down:
	migrate -path cmd/db/migration -database ${DB_ADDR} -verbose down

migrate_down1:
	migrate -path cmd/db/migration -database ${DB_ADDR} -verbose down 1

seed:
	@echo "Seeding database..."
	@go run cmd/db/seed/main.go
	@echo "Database seeded!"
