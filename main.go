package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"Webex.API.Integration.And.Visualization/api"
	"Webex.API.Integration.And.Visualization/persist"
)

func main() {
	// Initialize the sqlite db.
	db, err := sql.Open("sqlite3", "./persist/webex.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	p, err := persist.NewPersist(db)
	if err != nil {
		log.Fatal(err)
	}

	// Start the server and close if error occurs
	if err := api.WebexApplicationServer(p); err != nil {
		panic(err)
	}
}
