include .env
export

build:
	@go build -o bin/api/main cmd/main.go

run: build
	@./bin/api/main

migrate:
	@dotenv -f .env migrate -path db/migrations -database "$(DB_CONN_STR)" -verbose up

drop:
	@dotenv -f .env migrate -path db/migrations -database "$(DB_CONN_STR)" -verbose down

