postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=admin -d postgres:12-alpine

postgresbanknetwork:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=admin -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb  simple_bank

migrateup:
	migrate -path db/migration/ -database "postgresql://root:admin@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration/ -database "postgresql://root:admin@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration/ -database "postgresql://root:admin@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration/ -database "postgresql://root:admin@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mockdb:
	mockgen -package mockdb -destination db/mock/store.go simplebank/db/sqlc Store

dockerbuildimage:
	docker build -t simplebank:latest .

dockerrunimage:
	docker run --name simplebank --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:admin@postgres12:5432/simple_bank?sslmode=disable" simplebank:latest

dockercreateconnectnetwork:
	docker network connect bank-network postgres12

inspectbanknetwork:
	docker network inspect bank-network

.PHONY: postgres createdxb dropdb migrateup migratedown sqlc test server mockdb migrateup1 migratedown1 dockerbuildimage dockerrunimage dockercreateconnectnetwork inspectbanknetwork postgresbanknetwork