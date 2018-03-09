package main

import (
	"log"

	"github.com/tuommii/jumbo/database"
	"github.com/tuommii/jumbo/server"
)

func main() {
	db, err := database.NewSQLiteDB("jumbo.db")
	if err != nil {
		panic(err)
	}
	defer db.Connection.Close()

	server := server.Create(db)
	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
