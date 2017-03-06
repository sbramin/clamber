package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

var gdb *bolt.DB

func boltOn(c string) *bolt.DB {
	db, err := bolt.Open("wc.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	if c == "crawl" {

		err := db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(baseURL))
			if bucket != nil {
				tx.DeleteBucket([]byte(baseURL))
			}
			return nil
		})
		if err != nil {
			log.Print(err)
		}

		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(baseURL))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		})
		if err != nil {
			log.Print(err)
		}
	}

	gdb = db
	return db
}

func boltOff(db *bolt.DB) {
	db.Close()
}

func boltUp() []string {

	var pages []string
	err := gdb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(baseURL))

		if bucket == nil {
			fmt.Println("You haven't crawled that site")
		} else {
			c := bucket.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				pages = append(pages, string(v))
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)

	}
	return pages
}

func boltDown(page pageType) {
	buf, err := json.Marshal(page)
	if err != nil {
		log.Print(err)
	}
	err = gdb.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(baseURL))
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(page.URL), buf)
		if err != nil {
			log.Print(err)
		}
		return nil
	})
}
