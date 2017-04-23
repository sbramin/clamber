package main

import (
	"testing"

	"github.com/boltdb/bolt"
)

func TestNilBucket(t *testing.T) {
	b := "FakeBuket123"
	db, _ := bolt.Open("/tmp/clamber.db", 0600, nil)
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b))
		if bucket != nil {
			t.Error("Shouldn't be able to find a fake bucket")
		}
		return nil
	})
}

func TestBucketCreate(t *testing.T) {
	_, err := bolt.Open("/root/bolt.db", 0600, nil)
	if err == nil {
		t.Error("Are you running this as root?, should have failed with", err)
	}
}
