package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

type boltDB struct {
	*bolt.DB
}

func createDB(baseURL *string, job string) *boltDB {
	bdb, err := bolt.Open("clamber.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	if job == "crawl" {

		err := bdb.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(baseURL))
			if bucket != nil {
				tx.DeleteBucket([]byte(baseURL))
			}
			return nil
		})
		if err != nil {
			log.Print(err)
		}

		err = bdb.Update(func(tx *bolt.Tx) error {
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
	db := &boltDB{DB: bdb}
	return db
}

func (db *boltDB) Off() {
	db.Close()
}

func (db *boltDB) Read() []string {

	var pages []string
	err := db.View(func(tx *bolt.Tx) error {
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

func (db *boltDB) Write(page pageType) {
	buf, err := json.Marshal(page)
	if err != nil {
		log.Print(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
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
	if err != nil {
		log.Print(err)
	}

}
