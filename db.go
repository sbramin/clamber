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

// setupDB - creates a new boltDB if none existed and initializes a bucket for the site being
// worked on.  If a crawl operation is being performed and it has previously crawled that site
// it will first delete the previous bucket.
func setupDB() (*boltDB, error) {
	bdb, err := bolt.Open("clamber.db", 0600, nil)
	if err != nil {
		return &boltDB{}, fmt.Errorf("could not open db file - %s", err)
	}
	return &boltDB{DB: bdb}, nil
}

// CreateBucket is a wrapper around boltDBs Update/Create bucket methods that
// first removes a bucket if it already exists.
func (db *boltDB) CreateBucket(baseURL *string) (err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(*baseURL))
		if b != nil {
			err = tx.DeleteBucket([]byte(*baseURL))
		}
		if err != nil {
			return fmt.Errorf("could not clean out existing bucket %s", err)
		}

		b, err = tx.CreateBucketIfNotExists([]byte(*baseURL))
		if err != nil {
			return fmt.Errorf("create bucket %s", err)
		}
		return err
	})
	return err
}

// Reader method for boltDB type
func (db *boltDB) Reader(baseURL *string) []string {

	var ps []string
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(*baseURL))

		if bucket == nil {
			return fmt.Errorf("you haven't crawled that site")
		}
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			ps = append(ps, string(v))

		}
		return nil
	})
	if err != nil {
		log.Fatal(err)

	}
	return ps
}

// Writer method for boltDB type
func (db *boltDB) Writer(baseURL *string, p page) {
	buf, err := json.Marshal(p)
	if err != nil {
		log.Print(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(*baseURL))
		err = bucket.Put([]byte(p.URL), buf)
		return err
	})
	if err != nil {
		log.Print(err)
	}

}
