package main

import (
	"encoding/json"
	"fmt"
)

func review() {
	if pretty {
		prettyJSON()
	} else {
		plainJSON()
	}
}

func prettyJSON() {
	var page pageType

	for _, p := range boltUp() {
		json.Unmarshal([]byte(p), &page)
		j, _ := json.MarshalIndent(&page, "", " ")
		fmt.Println(string(j))
	}

}

func plainJSON() {
	for _, p := range boltUp() {
		fmt.Println(p)
	}

}
