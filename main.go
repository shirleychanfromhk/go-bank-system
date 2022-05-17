package main

import (
	"database/sql"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:admin@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {

	connection, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("DB Connection [ Failed ]: ", err)
	}

	store := db.NewStore(connection)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}

}
