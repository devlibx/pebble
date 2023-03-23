package common

import (
	"bufio"
	"fmt"
	aws2 "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cockroachdb/pebble/vfs"
	"os"
	"strings"
)

type s3HelperImpl struct {
	bucket     string
	filePrefix string
	*s3manager.Uploader
	*s3.S3
}

func (s *s3HelperImpl) skipS3Upload(name string) bool {
	if strings.HasSuffix(name, ".log") {
		return true
	} else if strings.HasSuffix(name, ".dbtmp") {
		return true
	}
	return false
}

func (s *s3HelperImpl) SyncFileToS3(file vfs.File, name string) error {
	if s.skipS3Upload(name) {
		return nil
	}
	out, err := s.Upload(&s3manager.UploadInput{
		Body:   bufio.NewReader(file),
		Bucket: aws2.String(s.bucket),
		Key:    aws2.String(s.filePrefix + "/" + name),
	})
	fmt.Println("Cloud file close: name=", name, out)
	return err
}

func (s *s3HelperImpl) DeleteS3File(name string) error {
	_, err := s.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws2.String(s.bucket),
		Key:    aws2.String(s.filePrefix + "/" + name),
	})
	return err
}

func NewS3Helper(options CloudFsOption) (*s3HelperImpl, error) {
	awsSession, _ := session.NewSession(&aws2.Config{
		Region: aws2.String("ap-south-1")},
	)
	s := &s3HelperImpl{
		bucket:     os.Getenv("S3_BUCKET"),
		filePrefix: options.BasePath,
		S3:         s3.New(awsSession),
		Uploader:   s3manager.NewUploader(awsSession),
	}
	return s, nil
}
