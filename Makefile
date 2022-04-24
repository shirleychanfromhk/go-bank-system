postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=admin -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb  simple_bank

migrateup:
	migrate -path /Users/cheung/dev/simplebank/db/migration/ -database "postgresql://root:admin@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path /Users/cheung/dev/simplebank/db/migration/ -database "postgresql://root:admin@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: postgres createdxb dropdb migrateup migratedown sqlc