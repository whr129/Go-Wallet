DB_URL=postgresql://root:secret@localhost:5432/transaction_service?sslmode=disable

network:
	docker network create go-wallet-network

postgres:
	docker run --name postgres17 --network go-wallet-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine

createdb:
	docker exec -it postgres17 createdb --username=root --owner=root transaction_service

dropdb:
	docker exec -it postgres17 dropdb transaction_service

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

server:
	go run server/main.go

dockerbuild:
	docker build -t transaction-service:latest .

dockerrun:
	docker run --name transaction-service --network go-wallet-network -p 8083:8080 -d transaction-service:latest

.PHONY: network postgres migrateup createdb dropdb migratedown migrateup1 migratedown1 new_migration sqlc server dockerbuild dockerrun	