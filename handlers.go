package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func clamberIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "")
}

func clamberShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := vars["url"]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	for _, p := range db.Reader(url) {
		fmt.Fprintln(w, p)
	}
	fmt.Fprintln(w, url)
}

func clamber(w http.ResponseWriter, r *http.Request) {
	var s site
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &s); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	db.CreateBucket(s.URL)
	goCrawl(db, s.URL)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(s); err != nil {
		panic(err)
	}

}
