package aws

import (
	"github.com/cockroachdb/pebble/cloud/common"
	"github.com/cockroachdb/pebble/vfs"
	"os"
	"strings"
)

type CloudFile struct {
	file     vfs.File
	name     string
	s3Helper S3Helper
	options  common.CloudFsOption
}

func NewCloudFile(baseFile vfs.File, name string, s3Helper S3Helper, options common.CloudFsOption) (vfs.File, error) {
	return &CloudFile{file: baseFile, name: name, s3Helper: s3Helper, options: options}, nil
}

func (c *CloudFile) Close() error {
	err := c.s3Helper.SyncFileToS3(c.file, c.name)
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
		_ = c.s3Helper.SyncFileToS3(c.file, c.name)
	}
	return c.file.Sync()
}

func (c *CloudFile) SyncTo(length int64) (fullSync bool, err error) {
	if strings.Contains(c.name, "MANIFEST") {
		_ = c.s3Helper.SyncFileToS3(c.file, c.name)
	}
	return c.file.SyncTo(length)
}

func (c *CloudFile) SyncData() error {
	if strings.Contains(c.name, "MANIFEST") {
		_ = c.s3Helper.SyncFileToS3(c.file, c.name)
	}
	return c.file.SyncData()
}

func (c *CloudFile) Prefetch(offset int64, length int64) error {
	return c.file.Prefetch(offset, length)
}

func (c *CloudFile) Fd() uintptr {
	return c.file.Fd()
}
