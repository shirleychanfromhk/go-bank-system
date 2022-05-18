package main

import (
	"database/sql"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/db/util"

	_ "github.com/lib/pq"
)

func main() {
	viberConfig, err := util.LoadViberConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	connection, err := sql.Open(viberConfig.DBDriver, viberConfig.DBSource)
	if err != nil {
		log.Fatal("DB Connection [ Failed ]: ", err)
	}

	store := db.NewStore(connection)
	server := api.NewServer(store)

	err = server.Start(viberConfig.ServerAddress)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}

}
