DB_URL=postgresql://root:secret@localhost:5432/auth_service?sslmode=disable

network:
	docker network create go-wallet-network

postgres:
	docker run --name postgres17 --network go-wallet-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine

createdb:
	docker exec -it postgres17 createdb --username=root --owner=root auth_service

dropdb:
	docker exec -it postgres17 dropdb auth_service

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

.PHONY: network postgres migrateup migratedown migrateup1 migratedown1 new_migration sqlc server