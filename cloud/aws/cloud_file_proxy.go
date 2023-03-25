package aws

import (
	"github.com/cockroachdb/pebble/cloud/common"
	"github.com/cockroachdb/pebble/vfs"
	"os"
)

type CloudFileProxy struct {
	name     string
	s3Helper S3Helper
	options  common.CloudFsOption
}

func NewCloudFileProxy(name string, s3Helper S3Helper, options common.CloudFsOption) (vfs.File, error) {
	return &CloudFileProxy{name: name, s3Helper: s3Helper, options: options}, nil
}

func (c CloudFileProxy) Close() error {
	return nil
}

func (c CloudFileProxy) Read(p []byte) (n int, err error) {
	panic("implement me")
}

func (c CloudFileProxy) ReadAt(p []byte, off int64) (n int, err error) {
	panic("implement me")
}

func (c CloudFileProxy) Write(p []byte) (n int, err error) {
	panic("implement me")
}

func (c CloudFileProxy) Preallocate(offset, length int64) error {
	panic("implement me")
}

func (c CloudFileProxy) Stat() (os.FileInfo, error) {
	panic("implement me")
}

func (c CloudFileProxy) Sync() error {
	panic("implement me")
}

func (c CloudFileProxy) SyncTo(length int64) (fullSync bool, err error) {
	panic("implement me")
}

func (c CloudFileProxy) SyncData() error {
	panic("implement me")
}

func (c CloudFileProxy) Prefetch(offset int64, length int64) error {
	panic("implement me")
}

func (c CloudFileProxy) Fd() uintptr {
	panic("implement me")
}
