package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// review calls either plain or pretty json for a previously crawled site form boltDB
func review(db *boltDB, baseURL *string, pretty *bool) {
	if *pretty {
		prettyJSON(db, baseURL)
	} else {
		plainJSON(db, baseURL)
	}
}

// prettyJSON retrieves records from boltdb of a previously crawled site and unmarshals
// the json inorder to prettily print it to the terminal.
func prettyJSON(db *boltDB, baseURL *string) {
	var page pageType

	for _, p := range db.Read(baseURL) {
		err := json.Unmarshal([]byte(p), &page)
		if err != nil {
			log.Print(err)
		}
		j, _ := json.MarshalIndent(&page, "", " ")
		fmt.Println(string(j))
	}
}

// plainJSON retrieves the json records for a previously crawled site and prints them
// to the terminal.
func plainJSON(db *boltDB, baseURL *string) {
	for _, p := range db.Read(baseURL) {
		fmt.Println(p)
	}
}
