package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cockroachdb/pebble/vfs"
	"io"
	"os"
)

type CloudFS struct {
	wrapperFs vfs.FS
	options   CloudFsOption
	s3Client  *s3.S3
}

type CloudFsOption struct {
	BasePath string
}

func (c *CloudFS) Create(name string) (vfs.File, error) {
	if f, err := c.wrapperFs.Create(name); err == nil {
		return NewCloudFile(f, name, c.options)
	} else {
		return nil, err
	}
}

func (c *CloudFS) Link(oldname, newname string) error {
	return c.wrapperFs.Link(oldname, newname)
}

func (c *CloudFS) Open(name string, opts ...vfs.OpenOption) (vfs.File, error) {
	return c.wrapperFs.Open(name, opts...)
}

func (c *CloudFS) OpenDir(name string) (vfs.File, error) {
	return c.wrapperFs.OpenDir(name)
}

func (c *CloudFS) Remove(name string) error {
	if out, err := c.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(c.options.BasePath + "/" + name),
	}); err == nil {
		fmt.Println("Delete S3 file", out)
	}
	return c.wrapperFs.Remove(name)
}

func (c *CloudFS) RemoveAll(name string) error {
	return c.wrapperFs.RemoveAll(name)
}

func (c *CloudFS) Rename(oldname, newname string) error {
	if baseFile, err := c.wrapperFs.Create(oldname); err == nil {
		if oldFile, err := NewCloudFile(baseFile, oldname, c.options); err == nil {
			(oldFile.(*CloudFile)).updateToS3(newname)
		}
	}
	return c.wrapperFs.Rename(oldname, newname)
}

func (c *CloudFS) ReuseForWrite(oldname, newname string) (vfs.File, error) {
	return c.wrapperFs.ReuseForWrite(oldname, newname)
}

func (c *CloudFS) MkdirAll(dir string, perm os.FileMode) error {
	return c.wrapperFs.MkdirAll(dir, perm)
}

func (c *CloudFS) Lock(name string) (io.Closer, error) {
	return c.wrapperFs.Lock(name)
}

func (c *CloudFS) List(dir string) ([]string, error) {
	return c.wrapperFs.List(dir)
}

func (c *CloudFS) Stat(name string) (os.FileInfo, error) {
	return c.wrapperFs.Stat(name)
}

func (c *CloudFS) PathBase(path string) string {
	return c.wrapperFs.PathBase(path)
}

func (c *CloudFS) PathJoin(elem ...string) string {
	return c.wrapperFs.PathJoin(elem...)
}

func (c *CloudFS) PathDir(path string) string {
	return c.wrapperFs.PathDir(path)
}

func (c *CloudFS) GetDiskUsage(path string) (vfs.DiskUsage, error) {
	return c.wrapperFs.GetDiskUsage(path)
}

func NewCloudFS(fs vfs.FS, options CloudFsOption) vfs.FS {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1")},
	)

	cfs := &CloudFS{
		wrapperFs: fs,
		options:   options,
		s3Client:  s3.New(sess),
	}
	return cfs
}
