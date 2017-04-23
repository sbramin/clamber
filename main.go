package main

import (
	"log"
	"net/http"
)

var db *boltDB

func main() {
	var err error
	db, err = setupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
