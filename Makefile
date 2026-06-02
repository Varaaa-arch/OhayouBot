.PHONY: dev up down migrate sqlc build

up:
	docker compose up --build

down:
	docker compose down

migrate:
	migrate -path apps/api/migrations -database "$$DB_URL" up

sqlc:
	cd apps/api && sqlc generate

build:
	docker compose build
