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
	// example cronjob which print hello world for every second. Uncomment the line to taste it.
	// go cronjob.StartCronJob("@every 1s", &cronjob.ExampleJob{}, time.Local)

	// runing cronjob at background at every 0000,0800,1600 UTC
	// go cronjob.StartCronJob("0 0,8,16 * * *", &cronjob.ExchangeRateUpdateJob{}, time.UTC)

	config, err := util.LoadViberConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	connection, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("DB Connection [ Failed ]: ", err)
	}

	store := db.NewStore(connection)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}

}
