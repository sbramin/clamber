package main

import (
	"encoding/json"
	"fmt"
)

func review(db *boltDB, baseURL *string, pretty *bool) {
	if *pretty {
		prettyJSON(db, baseURL)
	} else {
		plainJSON(db, baseURL)
	}
}

func prettyJSON(db *boltDB, baseURL *string) {
	var page pageType

	for _, p := range db.Read(baseURL) {
		json.Unmarshal([]byte(p), &page)
		j, _ := json.MarshalIndent(&page, "", " ")
		fmt.Println(string(j))
	}
}

func plainJSON(db *boltDB, baseURL *string) {
	for _, p := range db.Read(baseURL) {
		fmt.Println(p)
	}
}
