package main

import (
	"encoding/json"
	"fmt"
)

func review(pretty bool, db *boltDB) {
	if pretty {
		prettyJSON(db)
	} else {
		plainJSON(db)
	}
}

func prettyJSON(db *boltDB) {
	var page pageType

	for _, p := range db.Read() {
		json.Unmarshal([]byte(p), &page)
		j, _ := json.MarshalIndent(&page, "", " ")
		fmt.Println(string(j))
	}
}

func plainJSON(db *boltDB) {
	for _, p := range db.Read() {
		fmt.Println(p)
	}
}
