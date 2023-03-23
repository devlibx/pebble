package aws

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cockroachdb/pebble/vfs"
	"os"
	"strings"
)

type CloudFile struct {
	file vfs.File
	name string
	*s3manager.Uploader
	options CloudFsOption
}

func NewCloudFile(baseFile vfs.File, name string, options CloudFsOption) (vfs.File, error) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1")},
	)
	uploader := s3manager.NewUploader(sess)

	return &CloudFile{file: baseFile, name: name, Uploader: uploader, options: options}, nil
}

func (c *CloudFile) updateToS3(name string) error {
	if SkipS3Upload(name) {
		fmt.Println("Skipping file to S3: name=", name)
		return nil
	}
	out, err := c.Upload(&s3manager.UploadInput{
		Body:   bufio.NewReader(c.file),
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(c.options.BasePath + "/" + name),
	})
	fmt.Println("Cloud file close: name=", name, out)
	return err
}

func (c *CloudFile) Close() error {
	err := c.updateToS3(c.name)
	err = c.file.Close()
	return err
}

func (c *CloudFile) Read(p []byte) (n int, err error) {
	return c.file.Read(p)
}

func (c *CloudFile) ReadAt(p []byte, off int64) (n int, err error) {
	return c.file.ReadAt(p, off)
}

func (c *CloudFile) Write(p []byte) (n int, err error) {
	return c.file.Write(p)
}

func (c *CloudFile) Preallocate(offset, length int64) error {
	return c.file.Preallocate(offset, length)
}

func (c *CloudFile) Stat() (os.FileInfo, error) {
	return c.file.Stat()
}

func (c *CloudFile) Sync() error {
	if strings.Contains(c.name, "MANIFEST") {
		_ = c.updateToS3(c.name)
	}
	return c.file.Sync()
}

func (c *CloudFile) SyncTo(length int64) (fullSync bool, err error) {
	if strings.Contains(c.name, "MANIFEST") {
		_ = c.updateToS3(c.name)
	}
	return c.file.SyncTo(length)
}

func (c *CloudFile) SyncData() error {
	if strings.Contains(c.name, "MANIFEST") {
		_ = c.updateToS3(c.name)
	}
	return c.file.SyncData()
}

func (c *CloudFile) Prefetch(offset int64, length int64) error {
	return c.file.Prefetch(offset, length)
}

func (c *CloudFile) Fd() uintptr {
	return c.file.Fd()
}
