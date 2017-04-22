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

// createDB - creates a new boltDB if none existed and initializes a bucket for the site being
// worked on.  If a crawl operation is being performed and it has previously crawled that site
// it will first delete the previous bucket.
func createDB(baseURL, job *string) *boltDB {
	bdb, err := bolt.Open("clamber.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	if *job == "crawl" {

		err := bdb.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(*baseURL))
			if bucket != nil {
				err = tx.DeleteBucket([]byte(*baseURL))
			}
			if err != nil {
				return fmt.Errorf("could not delete bucket %s", err)
			}
			return nil
		})
		if err != nil {
			log.Print(err)
		}
	}

	err = bdb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(*baseURL))
		if err != nil {
			return fmt.Errorf("create bucket %s", err)
		}
		return nil
	})
	if err != nil {
		log.Print(err)
	}
	db := &boltDB{DB: bdb}
	return db
}

func (db *boltDB) Off() {
	err := db.Close()
	if err != nil {
		log.Print(err)
	}
}

// Read method for boltDB type
func (db *boltDB) Read(baseURL *string) []string {

	var pages []string
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(*baseURL))

		if bucket == nil {
			return fmt.Errorf("You haven't crawled that site")
		}
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			pages = append(pages, string(v))

		}
		return nil
	})
	if err != nil {
		log.Fatal(err)

	}
	return pages
}

// Write method for boltDB type
func (db *boltDB) Write(baseURL *string, page pageType) {
	buf, err := json.Marshal(page)
	if err != nil {
		log.Print(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(*baseURL))
		err = bucket.Put([]byte(page.URL), buf)
		return err
	})
	if err != nil {
		log.Print(err)
	}

}
