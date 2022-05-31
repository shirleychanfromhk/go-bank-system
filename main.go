package main

import (
	"database/sql"
	"log"
	"simplebank/api"
	"simplebank/cronjob"
	db "simplebank/db/sqlc"
	"simplebank/db/util"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// runing cronjob at background at every 0000,0800,1600 UTC
	go cronjob.StartCronJob("0 0,8,16 * * *", &cronjob.ExchangeRateUpdateJob{}, time.Local)

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
