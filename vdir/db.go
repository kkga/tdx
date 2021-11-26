package vdir

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
)

var bucketName = "todos"

func dbPath() (dbPath string, err error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("get cache dir: %w", err)
	}
	cacheDir := filepath.Join(userCacheDir, "tdx")
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return "", fmt.Errorf("create cache dir: %w", err)
	}
	dbPath = filepath.Join(cacheDir, "db")
	return
}

func updateDB() error {
	path, err := dbPath()
	if err != nil {
		return err
	}

	db, err := bolt.Open(path, 0644, nil)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	if err != nil {
		return err
	}

	key := []byte("hey")
	val := []byte("world")

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		err = bucket.Put(key, val)
		if err != nil {
			return err
		}
		return nil
	})

	return nil
}

func viewDB() (string, error) {
	path, err := dbPath()
	if err != nil {
		return "", err
	}

	db, err := bolt.Open(path, 0644, nil)
	if err != nil {
		return "", fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	key := []byte("hey")
	var val []byte

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found", bucketName)
		}
		val = bucket.Get(key)

		return nil
	})

	if err != nil {
		return "", err
	}
	fmt.Println(string(key), string(val))

	return string(val), nil
}
