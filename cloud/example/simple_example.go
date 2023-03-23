package main

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cockroachdb/pebble"
	aws2 "github.com/cockroachdb/pebble/cloud/aws"
	"github.com/cockroachdb/pebble/vfs"
	"log"
	"os"
	"strings"
)

func s3_main() {

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1")},
	)

	downloader := s3manager.NewUploader(sess)

	file, err := os.Open("/tmp/t/test.txt")
	if err != nil {
		panic(err)
	}

	out, err := downloader.Upload(&s3manager.UploadInput{
		Body:   bufio.NewReader(file),
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String("harish_1"),
	})
	println(err, out)

}
func main() {
	if true {
		// s3_main()
		// return
	}

	id := "5"

	baseFs := vfs.Default
	baseFs = aws2.NewCloudFS(baseFs, aws2.CloudFsOption{BasePath: "project_" + id})
	baseFs = vfs.WithLogging(baseFs, func(_fmt string, args ...interface{}) {
		if strings.Contains(_fmt, "sync-data") {
			return
		}
		fmt.Printf(_fmt+"\n", args...)
	})
	db, err := pebble.Open("/tmp/demo_"+id, &pebble.Options{
		// FS: pAws.NewCloudFS(baseFs),
		FS: baseFs,
	})
	if err != nil {
		log.Fatal(err)
	}

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
