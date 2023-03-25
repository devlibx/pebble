package main

import (
	"fmt"
	"github.com/cockroachdb/pebble"
	aws2 "github.com/cockroachdb/pebble/cloud/aws"
	"github.com/cockroachdb/pebble/cloud/common"
	"github.com/cockroachdb/pebble/vfs"
	"log"
	"strings"
	"time"
)

func main() {
	id := "12"

	baseFs := vfs.Default
	baseFs, err := aws2.NewCloudFS(baseFs, common.CloudFsOption{BasePath: "project_" + id})
	if err != nil {
		panic(err)
	}
	baseFs = vfs.WithLogging(baseFs, func(_fmt string, args ...interface{}) {
		if strings.Contains(_fmt, "sync-data") {
			return
		}
		// fmt.Printf(_fmt+"\n", args...)
	})
	db, err := pebble.Open("/tmp/demo_"+id, &pebble.Options{
		// FS: pAws.NewCloudFS(baseFs),
		FS: baseFs,
	})
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for i := 0; i < 10000000; i++ {
			key := []byte(fmt.Sprintf("hello_%d", i))
			if _, _, err := db.Get(key); err == nil {
				fmt.Println("Data for Key=", string(key))
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	key := []byte("")
	data := strings.Repeat("world", 10000)
	for i := 0; i < 10000000; i++ {
		key := []byte(fmt.Sprintf("hello_%d", i))
		if err := db.Set(key, []byte(data), pebble.Sync); err != nil {
			log.Fatal(err)
		}
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
