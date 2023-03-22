package main

import (
	"fmt"
	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/cloud/aws"
	"github.com/cockroachdb/pebble/vfs"
	"log"
)

func main() {
	baseFs := vfs.Default
	db, err := pebble.Open("demo", &pebble.Options{
		FS: aws.NewCloudFS(baseFs),
	})
	if err != nil {
		log.Fatal(err)
	}
	key := []byte("hello")
	if err := db.Set(key, []byte("world"), pebble.Sync); err != nil {
		log.Fatal(err)
	}
	value, closer, err := db.Get(key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s %s\n", key, value)
	if err := closer.Close(); err != nil {
		log.Fatal(err)
	}
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}
