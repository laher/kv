package kv

import (
	"fmt"

	"github.com/boltdb/bolt"
)

type KeyValue struct {
	db     *bolt.DB
	bucket string
}

func NewKeyValue(db *bolt.DB, bucket string) (*KeyValue, error) {
	return &KeyValue{db, bucket}, db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})
}

func (kv *KeyValue) Set(key, value string) error {
	return kv.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(kv.bucket))
		if bucket == nil {
			return fmt.Errorf("Bucket %s does not exist", kv.bucket)
		}
		return bucket.Put([]byte(key), []byte(value))
	})
}

func (kv *KeyValue) Get(key string) (string, error) {
	val := ""
	if err := kv.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(kv.bucket))
		if bucket == nil {
			return fmt.Errorf("Bucket %s does not exist", kv.bucket)
		}
		v := bucket.Get([]byte(key))
		if v == nil {
			return fmt.Errorf("Key %s does not exist", key)
		}
		val = string(v)
		return nil
	}); err != nil {
		return "", err
	}
	return val, nil
}

func (kv *KeyValue) Del(key string) error {
	return kv.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(kv.bucket))
		if bucket == nil {
			return fmt.Errorf("Bucket %s does not exist", kv.bucket)
		}
		return bucket.Delete([]byte(key))
	})
}
